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

// GetAccount finds a value using key in redis
func (s *RedisClient) GetAccount(user string) (*models.AccountStruct, error) {
	val, err := s.database.Get(context.Background(), user).Result()
	if err != nil {
		return &models.AccountStruct{}, errors.New("Account with that name does not exist.")
	}
	return &models.AccountStruct{
		User: user,
		Hash: val,
	}, nil
}

// Add account to redis
func (s *RedisClient) AddAccount(user, passhash string) (bool, error) {
	res, err := s.AccountExists(user)
	if err == nil || res {
		return false, errors.New("Account exists with that name.")
	}
	// add the new account
	err = s.database.Set(context.Background(), user, passhash, 0).Err()
	if err != nil {
		return false, errors.New("Failed to add account.")
	}
	return true, nil
}

// DeleteAccount removes an account from redis
func (s *RedisClient) DeleteAccount(user string) (bool, error) {
	_, err := s.database.Del(context.Background(), user).Result()
	if err != nil {
		return false, errors.New("failed to remove account")
	}
	return true, nil
}

// AccountExists checks if an account key exists in redis
func (s *RedisClient) AccountExists(user string) (bool, error) {
	res, err := s.database.Exists(context.Background(), user).Result()
	if err != nil || res == 0 {
		return false, errors.New("account doesnt exists")
	}
	return true, nil
}
