package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mshayler/accountapi/internal/models"
	"github.com/mshayler/accountapi/internal/service"
	"time"
)

// Endpoints for gokit struct
type Endpoints struct {
	CreateAccount endpoint.Endpoint
	Verify        endpoint.Endpoint
	Login         endpoint.Endpoint
}

// MakeEndpoints returns a new instance of Endpoints.
func MakeEndpoints(svc service.Service, timeout time.Duration) Endpoints {
	return Endpoints{
		CreateAccount: makeCreateAccountEndpoint(svc),
		Verify:        makeVerifyEndpoint(svc),
		Login:         makeLoginEndpoint(svc),
	}
}

func makeCreateAccountEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(models.AccountRequest)
		res, err := svc.CreateAccount(context.Background(), req)
		if err != nil {
			return models.AccountResponse{"Can't create account"}, err
		}
		return models.AccountResponse{res.Result}, nil
	}
}

func makeVerifyEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(models.VerifyRequest)
		res, err := svc.Verify(context.Background(), req)
		if err != nil {
			return models.VerifyResponse{Result: "Can't verify account :("}, err
		}
		return models.VerifyResponse{Result: res.Result}, nil
	}
}
func makeLoginEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(models.AccountRequest)
		res, err := svc.Login(context.Background(), req)
		if err != nil {
			return models.LoginResponse{""}, err
		}
		return models.LoginResponse{res.Token}, nil
	}
}
