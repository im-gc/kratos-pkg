package redis

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/durationpb"
)

// Nil reply returned by Redis when key does not exist.
const Nil = redis.Nil

type (
	Redis struct {
		rdb *redis.Client
	}
	Option interface {
		GetNetwork() string
		GetAddr() string
		GetPassword() string
		GetSelectDb() int32
		GetReadTimeout() *durationpb.Duration
		GetWriteTimeout() *durationpb.Duration
		GetMaxIdleConns() int32
		GetMinIdleConns() int32
		GetConnMaxIdleTime() *durationpb.Duration
		GetConnMaxLiftTime() *durationpb.Duration
	}
)

func NewRedis(opt Option, logger log.Logger) (*Redis, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            opt.GetAddr(),
		Password:        opt.GetPassword(),
		DB:              int(opt.GetSelectDb()),
		MaxIdleConns:    int(opt.GetMaxIdleConns()),
		MinIdleConns:    int(opt.GetMinIdleConns()),
		ConnMaxIdleTime: opt.GetConnMaxIdleTime().AsDuration(),
		ConnMaxLifetime: opt.GetConnMaxLiftTime().AsDuration(),
		ReadTimeout:     opt.GetReadTimeout().AsDuration(),
		WriteTimeout:    opt.GetWriteTimeout().AsDuration(),
	})
	return &Redis{rdb: rdb}, func() {}, nil
}

func (r *Redis) RDB(ctx context.Context) *redis.Client {
	return r.rdb
}

func ExistKey(err error) bool {
	if nil == err {
		return true
	}
	return !errors.Is(err, Nil)
}

func Error(err error) error {
	if yes := errors.Is(err, Nil); yes {
		return nil
	}
	return err
}
