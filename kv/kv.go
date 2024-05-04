package kv

import (
	"fmt"
	"time"
)

const (
	KVErrorDefault          = 10000
	KVErrorNotFound         = 10001
	KVErrorKeyNotProvider   = 10002
	KVErrorGlobalNotProvide = 10003
	KVErrorCacheDefense     = 10004
)

var (
	errKeyNotFound          = "KV缓存器[%s], 数据源查询失败,缓存key[%s]没命中"
	errFn4KeyNotProvider    = "KV缓存器[%s], 按key源数据查询器没有提供"
	errFn4GlobalNotProvider = "KV缓存器[%s], 全局源数据查询器没有提供"
	errCacheDefense         = "KV缓存器[%s], 查询源数据异常,本地缓存兜底.异常原因:%s"
	errDefault              = "KV缓存器[%s]异常: %s"

	ErrNotInit                = NewErrDefault("", "KV is not instantiated")
	ErrDataQueryerNotProvider = NewErrFn4GlobalNotProvider("the kv querier does not provided")

	DefaultKey   = "_kv_default_"
	NeverExpired = int64(-1)
)

type (
	Cancel                           func() bool
	FetchGlobal[K comparable, V any] func() (map[K]V, error)
	FetchKey[K comparable, V any]    func(K) (map[K]V, error)
	DataCheck[K comparable, V any]   func(map[K]V) error
	MetricReport                     func(kvEname, kvName string, err error)

	KV[K comparable, V any] interface {
		// 直接获取缓存
		GetCacheByKey(K) (V, KVError)

		// 获取缓存，如果没有命中则回源
		GetCacheOrElseSourceByKey(K) (V, KVError)

		// 直接获取缓存
		GetCache() (map[K]V, KVError)

		// 获取缓存，如果没有命中则回源
		GetCacheOrElseSource() (map[K]V, KVError)

		// 直接回源
		GetSourceByKey(K) (V, KVError)

		// 直接回源
		GetSource() (map[K]V, KVError)

		// 查询数据源，如果查询失败 则通过本地缓存兜底
		// 如果查询源数据失败，则走本地缓存。 error将返回ErrCacheDefense
		GetSourceOrElseCache() (map[K]V, KVError)

		// 查询数据源，如果查询失败 则通过本地缓存兜底
		// 如果查询源数据失败，则走本地缓存。 error将返回ErrCacheDefense
		GetSourceOrElseCacheByKey(K) (V, KVError)
	}

	KVError interface {
		error
		Code() int64
	}
	ErrKeyNotFound struct {
		code int64
		msg  string
	}
	ErrFn4KeyNotProvider struct {
		code int64
		msg  string
	}
	ErrFn4GlobalNotProvider struct {
		code int64
		msg  string
	}
	ErrCacheDefense struct {
		code int64
		msg  string
	}
	ErrDefault struct {
		code int64
		msg  string
	}
)

func NewErrDefault(name, msg string) *ErrDefault {
	return &ErrDefault{
		code: KVErrorDefault,
		msg:  fmt.Sprintf(errDefault, name, msg),
	}
}

func NewErrKeyNotFound(name string, key any) *ErrKeyNotFound {
	return &ErrKeyNotFound{
		code: KVErrorNotFound,
		msg:  fmt.Sprintf(errKeyNotFound, name, key),
	}
}

func NewErrFn4KeyNotProvider(name string) *ErrFn4KeyNotProvider {
	return &ErrFn4KeyNotProvider{
		code: KVErrorKeyNotProvider,
		msg:  fmt.Sprintf(errFn4KeyNotProvider, name),
	}
}

func NewErrFn4GlobalNotProvider(name string) *ErrFn4GlobalNotProvider {
	return &ErrFn4GlobalNotProvider{
		code: KVErrorGlobalNotProvide,
		msg:  fmt.Sprintf(errFn4GlobalNotProvider, name),
	}
}

func NewErrCacheDefense(name, msg string) *ErrCacheDefense {
	return &ErrCacheDefense{
		code: KVErrorCacheDefense,
		msg:  fmt.Sprintf(errCacheDefense, name, msg),
	}
}

func (e *ErrKeyNotFound) Code() int64 {
	return e.code
}

func (e *ErrKeyNotFound) Error() string {
	return e.msg
}

func (e *ErrFn4GlobalNotProvider) Code() int64 {
	return e.code
}

func (e *ErrFn4GlobalNotProvider) Error() string {
	return e.msg
}

func (e *ErrFn4KeyNotProvider) Code() int64 {
	return e.code
}

func (e *ErrFn4KeyNotProvider) Error() string {
	return e.msg
}

func (e *ErrCacheDefense) Code() int64 {
	return e.code
}

func (e *ErrCacheDefense) Error() string {
	return e.msg
}

func (e *ErrDefault) Code() int64 {
	return e.code
}

func (e *ErrDefault) Error() string {
	return e.msg
}

// ticker 定时器
func ticker(interval int, fn func(), cancel Cancel) {

	if interval <= 0 || (cancel != nil && cancel()) {
		return
	}

	fn()
	go func() {
		eventsTick := time.NewTicker(time.Duration(interval) * time.Second)
		defer eventsTick.Stop()
		for {
			if cancel != nil && cancel() {
				break
			}
			select {
			case <-eventsTick.C:
				fn()
			}
		}
	}()
}
