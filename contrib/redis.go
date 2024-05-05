package contrib

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	ErrRedisConnectionTimeout = errors.New("ERR_REDIS_CONNECTION_TIMEOUT")
)

type RedisConf interface {
	GetAddr() string
	GetPassword() string
	GetSelectDb() int32
}

// NewRedis ... 创建一个 Redis 实例
func NewRedis(c RedisConf, ctxlog *log.Helper) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.GetAddr(),
		Password: c.GetPassword(),
		DB:       int(c.GetSelectDb()),
	})

	// first connect timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		ctxlog.Error("connect redis timeout", "error", err)
		return nil, ErrRedisConnectionTimeout
	}

	return client, nil
}
