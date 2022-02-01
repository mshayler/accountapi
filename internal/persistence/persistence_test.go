package persistence

import "testing"

func TestRedisClient_DeleteAccount(t *testing.T) {
	client, err := NewRedisClient()
	if err != nil {
		panic("failed to create redis client")
	}

	_, err = client.DeleteAccount("test")
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

	res, err := client.AddAccount("test", testSha)

	if err != nil || res == false {
		panic("failed to add account")
	}

	res, err = client.AccountExists("test")
	if !res {
		panic("account doesnt exist.")
	}

	res, err = client.DeleteAccount("test")
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

	res, err := client.AddAccount("test", testSha)

	if err != nil || res == false {
		panic("failed to add account")
	}

	accInfo, err := client.GetAccount("test")
	if err != nil {
		panic("failed to get account")
	}

	if accInfo.Hash != testSha {
		panic("hashes do not match")
	}

	res, err = client.DeleteAccount("test")
	if err != nil {
		panic("failed to delete account")
	}

}
