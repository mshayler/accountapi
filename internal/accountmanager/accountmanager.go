package accountmanager

import (
	"github.com/mshayler/accountapi/internal/models"
	"github.com/mshayler/accountapi/internal/persistence"
	"github.com/op/go-logging"
)

// AccountManager is here to manage overarching account management, and login tokens for each user
type AccountManager interface {
	CreateAccount() func(user, pass string) (bool, error)
	LoginAccount() func(user, pass string) (string, error)
	VerifyAccount() func(user, token string) (bool, error)
}

type Manager struct {
	LoginMap    map[string]models.LoginStruct
	Logger      *logging.Logger
	Persistence persistence.Persistence
	AccountManager
}
