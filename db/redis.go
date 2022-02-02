package db

import (
	"github.com/go-redis/redis/v8"
)

func NewRedis() redis.Cmdable {
	// config := config.Get().RedisConfig
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	return rdb
}
