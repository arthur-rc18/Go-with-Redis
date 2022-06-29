package database

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

var Context = context.Background()

func ConnectRedis() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

}

func DatabaseConnection() *redis.Client {
	if db == nil {
		snc.Do(func() {
			db = ConnectRedis()
		})
	}
	return db
}

var (
	db *redis.Client

	snc sync.Once
)
