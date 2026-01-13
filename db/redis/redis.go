package db_redis

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once sync.Once
	rdb  *redis.Client
)

func Init() {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     "118.31.238.252:6379",
			Password: "redis_QPFnYb",
			DB:       0,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
	})
}

func Client() *redis.Client {
	if rdb == nil {
		panic("redis.Init not called")
	}
	return rdb
}
