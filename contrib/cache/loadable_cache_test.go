package cache_test

import (
	"context"
	"fmt"
	"github.com/imkouga/kratos-pkg/contrib/cache"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestBaseCache(t *testing.T) {

	//var c = 1

	cc := cache.New[int64, string](
		cache.WithTracing[int64, string](otel.GetTracerProvider()),
		cache.WithSize[int64, string](100),
		cache.WithExpiration[int64, string](2*time.Second), // 每5秒刷新一次缓存
		cache.WithRefreshAfterWrite[int64, string](func() map[int64]string {
			t.Logf("auto refresh kv")
			key := time.Now().Unix()

			//c++
			//if c > 3 {
			//	return nil
			//}

			return map[int64]string{
				key - 1: "新的value1",
				key:     "新的value2",
			}
		}),
	)

	// 每秒取一次缓存
	ticker := time.NewTicker(time.Microsecond * 100)
	for {
		select {
		case <-ticker.C:
			// do something
			kv := cc.GetALL(context.Background())
			//for k, v := range kv {
			//	t.Logf("k: %d\t v: %s", k, v)
			//}
			//t.Logf("--------------")
			if len(kv) <= 0 {
				t.Logf("kv is empty")
			}
			if len(kv) != 2 {
				panic("kv is not two-size.")
			}
			fmt.Println("ticker.")
		}
	}
}
