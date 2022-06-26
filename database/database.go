package database

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	db *redis.Client

	once sync.Once
)

var CTX = context.Background()

func ConnectRedis() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

}

func ConnectWithDB() *redis.Client {
	if db == nil {
		once.Do(func() {
			db = ConnectRedis()
		})
	}
	return db
}
