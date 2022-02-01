package persistence

import (
	"context"
	"testing"
)

var (
	ctx = context.Background()
)

func TestRedisClient_DeleteAccount(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		panic("failed to create redis client")
	}

	_, err = client.DeleteAccount(ctx, "test")
	if err != nil {
		panic("failed to delete account")
	}

}

func TestAccountExists(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		panic("failed to create redis client")
	}

	testSha := "5E884898DA28047151D0E56F8DC6292773603D0D6AABBDD62A11EF721D1542D8"

	res, err := client.AddAccount(ctx, "test", testSha)

	if err != nil || res == false {
		panic("failed to add account")
	}

	res, err = client.AccountExists(ctx, "test")
	if !res {
		panic("account doesnt exist.")
	}

	res, err = client.DeleteAccount(ctx, "test")
	if err != nil {
		panic("failed to delete account")
	}
}
func TestRedis(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		panic("failed to create redis client")
	}

	testSha := "5E884898DA28047151D0E56F8DC6292773603D0D6AABBDD62A11EF721D1542D8"

	res, err := client.AddAccount(ctx, "test", testSha)

	if err != nil || res == false {
		panic("failed to add account")
	}

	accInfo, err := client.GetAccount(ctx, "test")
	if err != nil {
		panic("failed to get account")
	}

	if accInfo.Hash != testSha {
		panic("hashes do not match")
	}

	res, err = client.DeleteAccount(ctx, "test")
	if err != nil {
		panic("failed to delete account")
	}

}
