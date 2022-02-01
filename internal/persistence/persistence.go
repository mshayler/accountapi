package persistence

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mshayler/accountapi/internal/models"
	"github.com/pkg/errors"
)


type RedisClient struct {
	database *redis.Client
}

// NewRedisClient generate a new client with database connection
func NewRedisClient() (*RedisClient, error) {
	// Should move this to a config and environment vars
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root", // no real password
		DB:       0,      // use default DB
	})

	// Ping to verify the client is working
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.New("Failed to create redis client")
	}

	return &RedisClient{
		rdb,
	}, nil
}