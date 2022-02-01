package service

import (
	"context"
	"github.com/mshayler/accountapi/internal/models"
)

// Service is the interface for functions
type Service interface {
	Login(ctx context.Context, req models.AccountRequest) (*models.LoginResponse, error)
	Verify(ctx context.Context, req models.VerifyRequest) (*models.VerifyResponse, error)
	CreateAccount(ctx context.Context, req models.AccountRequest) (*models.AccountResponse, error)
}
