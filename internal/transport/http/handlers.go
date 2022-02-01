package httpTransport

import (
	"context"
	"encoding/json"
	gkhttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/mshayler/accountapi/internal/endpoints"
	"github.com/mshayler/accountapi/internal/models"
	"github.com/op/go-logging"
	"net/http"
)

const (
	loginPath         = "/login"
	createAccountPath = "/create"
	verifyAccountPath = "/verify"
)

// Handlers is the HTTP handlers object.
type Handlers struct {
	http.Handler
	lgr *logging.Logger
}

// NewHandlers returns a new instance of HTTP handlers.
func NewHandlers(svcEndpoints endpoints.Endpoints, options []gkhttp.ServerOption, logger *logging.Logger) (*Handlers, error) {
	var handlers Handlers
	router := mux.NewRouter()

	router.Methods(http.MethodPost).Path(createAccountPath).Handler(gkhttp.NewServer(
		svcEndpoints.CreateAccount,
		handlers.decodeCreateAccountRequest,
		handlers.encodeResponse,
		options...,
	))

	router.Methods(http.MethodPost).Path(loginPath).Handler(gkhttp.NewServer(
		svcEndpoints.Login,
		handlers.decodeLoginRequest,
		handlers.encodeResponse,
		options...,
	))

	router.Methods(http.MethodGet).Path(verifyAccountPath).Handler(gkhttp.NewServer(
		svcEndpoints.Verify,
		handlers.decodeVerifyRequest,
		handlers.encodeResponse,
		options...,
	))
	handlers.Handler = router
	handlers.lgr = logger
	return &handlers, nil
}

// Decode Functions for Requests
func (h *Handlers) decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.lgr.Warning("Failed to decode account request")
		return nil, err
	}
	return request, nil
}

// decode verify
func (h *Handlers) decodeVerifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.lgr.Warning("Failed to decode verify request")
		return nil, err
	}
	return request, nil
}

// decode login
func (h *Handlers) decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request models.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.lgr.Warning("Failed to decode login request")
		return nil, err
	}
	return request, nil
}

// Default Encode Response
func (h *Handlers) encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
