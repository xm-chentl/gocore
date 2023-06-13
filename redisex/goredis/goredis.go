package goredis

import (
	"context"
	"time"

	"github.com/xm-chentl/gocore/redisex"

	"github.com/redis/go-redis/v9"
)

type redisImp struct {
	rdb *redis.Client
}

func (r redisImp) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r redisImp) Get(ctx context.Context, key string) (value string, err error) {
	cmd := r.rdb.Get(ctx, key)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}

	value, err = cmd.Result()
	return
}

func New(o redisex.Option) redisex.IRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     o.Addr,
		Password: o.Password,
	})

	return &redisImp{
		rdb: rdb,
	}
}
