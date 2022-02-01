package persistence

import "testing"

func TestRedis(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		panic("failed to create redis client")
	}

}
