package accountmanager

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"github.com/mshayler/accountapi/internal/models"
	"github.com/mshayler/accountapi/internal/persistence"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"time"
)

const (
	authKey = "trulioo"
)

// AccountManager is here to manage overarching account management, and login tokens for each user
type AccountManager interface {
	CreateAccount() func(ctx context.Context, user, pass string) (bool, error)
	LoginAccount() func(ctx context.Context, user, pass string) (string, error)
	VerifyAccount() func(ctx context.Context, user, token string) (bool, error)
}

type Manager struct {
	LoginMap    map[string]models.LoginStruct
	Logger      *logging.Logger
	Persistence persistence.Persistence
	AccountManager
}

// Handling for New Accounts Manager
func NewManager(lgr *logging.Logger) (*Manager, error) {
	// Check for logger pass
	if lgr == nil {
		return nil, errors.New("No provided logger")
	}
	// Check for rdb
	rdb, err := persistence.NewRedisClient()
	if err != nil {
		return nil, err
	}

	return &Manager{
		LoginMap:    make(map[string]models.LoginStruct),
		Logger:      lgr,
		Persistence: rdb,
	}, nil

}

// Create a new account with user and pass
func (c *Manager) CreateAccount(ctx context.Context, user, pass string) (bool, error) {
	// Sanitize input
	if user == "" || pass == "" {
		return false, errors.New("Missing parameter to create account.")
	}

	// Check if the account exists before we try and duplicate
	res, _ := c.Persistence.AccountExists(ctx, user)
	if res {
		c.Logger.Error("Tried to create duplicate account.")
		return false, errors.New("Invalid Account Name.")
	}

	// Create sha256 from pass
	phash := generateHash(pass)

	// Update the account information
	_, err := c.Persistence.AddAccount(ctx, user, phash)
	if err != nil {
		c.Logger.Error(err)
		return false, errors.New("Failed to add account.")
	}

	c.Logger.Info("Created Account for: ", user)
	return true, nil
}

// Login the account if user and pass are correct.
func (c *Manager) LoginAccount(ctx context.Context, user, pass string) (string, error) {
	// Validate input
	if user == "" || pass == "" {
		return "", errors.New("Missing parameters")
	}

	// Check the Account Exists
	res, err := c.Persistence.GetAccount(ctx, user)
	if err != nil {
		c.Logger.Error(err)
		return "", errors.New("Invalid credentials")
	}

	// Verify the supplied password is correct
	if res.Hash != generateHash(pass) {
		c.Logger.Error(err)
		return "", errors.New("Invalid credentials")
	}

	// Generate a Primitive Token and Save it to a Cache
	tkn, t := generateToken(user, res.Hash)
	c.LoginMap[user] = models.LoginStruct{
		Token:     tkn,
		Timestamp: t,
	}

	c.Logger.Info(user + " successfully logged in at: " + t)

	// Successful login
	return tkn, nil
}

// Verify that the account is authenticated providing user and token
// This should be an authenticated path that I validate bearer tokens, for sake of time I just made another route
// to verify the token is valid
func (c *Manager) VerifyAccount(ctx context.Context, user, token string) (bool, error) {
	// Validate input
	if user == "" || token == "" {
		return false, errors.New("Missing Parameters")
	}

	// Validate account exists
	res, err := c.Persistence.AccountExists(ctx, user)
	if err != nil || !res {
		c.Logger.Error(err)
		return false, errors.New("Could not verify account")
	}

	// Validate the token is equal and user has logged in
	tkn, ok := c.LoginMap[user]
	if !ok {
		return false, errors.New("Need to login first")
	}
	if tkn.Token != token {
		return false, errors.New("Invalid token supplied")
	}

	// TODO:: Check for duration to actually log the user out after some time.

	// Successful verification
	return true, nil
}

// Generate sha hash
func generateHash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// Primitive API Token Generation
func generateToken(user, phash string) (string, string) {
	t := time.Now().String()
	pref := generateHash(user + t)
	suff := generateHash(phash + authKey)
	tkn := pref + "." + suff
	return tkn, t
}
