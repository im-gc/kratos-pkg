package cache

import (
	"context"
	"github.com/szyhf/go-gcache/v2"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"
	"sync"
	"time"
)

type LoadableCache[K comparable, V any] interface {
	// Get a value pair to the cache data by key.
	Get(ctx context.Context, k K) (V, error)
	// GetALL returns all key-value pairs in the cache data.
	GetALL(context.Context) map[K]V
	// Values returns all values in the cache data.
	Values(context.Context) []V
	// Set a value pair to the cache data by key.
	Set(ctx context.Context, k K, v V) error
	// Purge clears all cache data.
	Purge(context.Context)
	// TryPurgeAndReload try to refresh cache data, if refresh result is nil, return false.
	TryPurgeAndReload(context.Context) bool
}

// loadableCache is a cache that can be refreshed.
type loadableCache[K comparable, V any] struct {
	mu sync.RWMutex

	tracer  oteltrace.Tracer   // 链路 provider
	traced  bool               // 是否将普通函数包装为 tracer
	c       gcache.Cache[K, V] // gcache 对象
	exp     time.Duration      // key 过期时间
	size    int                // 缓存大小,超出的缓存会被 evict
	block   bool               // 是否阻塞当前调用链
	refresh func() map[K]V     // 刷新缓存数据的函数
	ticker  *time.Ticker       // 定时器(用于过期刷缓存)
}

type Option[K comparable, V any] func(*loadableCache[K, V])

// WithRefreshAfterWrite refresh data provider
func WithRefreshAfterWrite[K comparable, V any](f func() map[K]V) Option[K, V] {
	return func(cb *loadableCache[K, V]) {
		cb.refresh = f
	}
}

// WithExpiration cache expiration, Automatically reload if timeout
func WithExpiration[K comparable, V any](exp time.Duration) Option[K, V] {
	return func(cb *loadableCache[K, V]) {
		cb.exp = exp
	}
}

// WithSize cache size limit.
func WithSize[K comparable, V any](size int) Option[K, V] {
	return func(cb *loadableCache[K, V]) {
		cb.size = size
	}
}

// WithBlock block call first RefreshAfterWrite.
func WithBlock[K comparable, V any]() Option[K, V] {
	return func(cb *loadableCache[K, V]) {
		cb.block = true
	}
}

// WithTracing enable otel tracing
func WithTracing[K comparable, V any](provider oteltrace.TracerProvider) Option[K, V] {
	return func(cb *loadableCache[K, V]) {
		cb.traced = true
		if provider == nil {
			provider = otel.GetTracerProvider()
		}
		cb.tracer = provider.Tracer("LoadableCache")
	}
}

func New[K comparable, V any](opts ...Option[K, V]) LoadableCache[K, V] {
	cache := &loadableCache[K, V]{
		exp:    10 * time.Second,
		size:   100,
		block:  false,
		traced: false,
	}
	// bind options
	for _, opt := range opts {
		opt(cache)
	}

	// 创建 gcache 对象
	cache.c = gcache.New[K, V](cache.size).
		Build()

	if cache.refresh != nil {
		cache.ticker = time.NewTicker(cache.exp)

		firstLoad := func() {
			if cache.refresh != nil {
				ret := cache.refresh()
				cache.putAll(ret)
			}
		}
		// block 状态不使用 goroutine, 卡住当前调用链等待结束
		if cache.block {
			firstLoad()
		} else {
			go firstLoad()
		}

		go cache.rf()
	}

	// 如果开启了链路追踪, 则包装函数
	if cache.traced {
		return &tracedWrapperLoadableCache[K, V]{
			cache,
		}
	}

	return cache
}

// Get a value pair to the cache data by key.
func (cb *loadableCache[K, V]) Get(ctx context.Context, k K) (V, error) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	v, err := cb.c.Get(k)
	return v, err
}

// GetALL returns all key-value pairs in the cache data.
func (cb *loadableCache[K, V]) GetALL(ctx context.Context) map[K]V {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return cb.c.GetALL(false)
}

// Values returns all values in the cache data.
func (cb *loadableCache[K, V]) Values(ctx context.Context) []V {
	m := cb.c.GetALL(false)
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

func (cb *loadableCache[K, V]) Purge(ctx context.Context) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.c.Purge()
}

func (cb *loadableCache[K, V]) TryPurgeAndReload(ctx context.Context) bool {
	defer func() {
		_ = recover()
	}()

	ret := cb.refresh()
	return cb.putAll(ret)
}

func (cb *loadableCache[K, V]) Set(ctx context.Context, k K, v V) error {
	return cb.c.Set(k, v)
}

func (cb *loadableCache[K, V]) rf() (V, bool) {
	for {
		select {
		case <-cb.ticker.C:
			if cb.refresh != nil {
				ret := cb.refresh()
				cb.putAll(ret)
			}
		}
	}
}

func (cb *loadableCache[K, V]) putAll(ret map[K]V) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// if ret len is zero, keep cache data
	// Tips: if you want to clear cache data, you can use cb.Purge()
	if len(ret) <= 0 {
		return false
	}

	cb.c.Purge()

	for k, v := range ret {
		_ = cb.c.Set(k, v)
	}

	return true
}
