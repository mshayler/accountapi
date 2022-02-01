package endpoints

import "github.com/go-kit/kit/endpoint"

// Endpoints for gokit struct
type Endpoints struct {
	CreateAccount endpoint.Endpoint
	Verify        endpoint.Endpoint
	Login         endpoint.Endpoint
}
