package redisex

import (
	"context"
	"time"
)

type IRedis interface {
	Set(context.Context, string, string, time.Duration) error
	Get(context.Context, string) (string, error)
}

type Option struct {
	Addr     string
	Password string
	UserName string
}
