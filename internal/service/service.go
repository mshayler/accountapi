package service

import (
	"context"
	"github.com/mshayler/accountapi/internal/accountmanager"
	"github.com/mshayler/accountapi/internal/models"
	"github.com/op/go-logging"
)

// Service is the interface for functions
type Service interface {
	Login(ctx context.Context, req models.AccountRequest) (*models.LoginResponse, error)
	Verify(ctx context.Context, req models.VerifyRequest) (*models.VerifyResponse, error)
	CreateAccount(ctx context.Context, req models.AccountRequest) (*models.AccountResponse, error)
}

type service struct {
	lgr    *logging.Logger
	accMan *accountmanager.Manager
}

// New returns a new instance of accounts service
func New(lgr *logging.Logger, acc *accountmanager.Manager) (Service, error) {
	return &service{
		lgr:    lgr,
		accMan: acc,
	}, nil
}

func (s *service) CreateAccount(ctx context.Context, req models.AccountRequest) (*models.AccountResponse, error) {
	s.lgr.Info("Recieved CreateAccount Request...")
	_, err := s.accMan.CreateAccount(req.User, req.Pass)
	if err != nil {
		return &models.AccountResponse{Result: err.Error()}, err
	}
	resp := &models.AccountResponse{Result: "Successfully Genereted Account for: " + req.User}
	return resp, nil
}

func (s *service) Login(ctx context.Context, req models.AccountRequest) (*models.LoginResponse, error) {
	s.lgr.Info("Recieved Login Request...")
	tkn, err := s.accMan.LoginAccount(req.User, req.Pass)
	if err != nil {
		return &models.LoginResponse{""}, err
	}
	resp := &models.LoginResponse{Token: tkn}
	return resp, nil
}
func (s *service) Verify(ctx context.Context, req models.VerifyRequest) (*models.VerifyResponse, error) {
	s.lgr.Info("Recieved Verify Request...")
	_, err := s.accMan.VerifyAccount(req.User, req.Token)
	if err != nil {
		return &models.VerifyResponse{Result: "Failed to Verify"}, err
	}
	resp := &models.VerifyResponse{Result: "Verified Account with Valid Credentials! Hello World!"}
	return resp, nil
}
