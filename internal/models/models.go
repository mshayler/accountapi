package models

// Model for account info
type AccountStruct struct {
	User string `json:"user"`
	Hash string `json:"hash"`
}

// Model for logged in device
type LoginStruct struct {
	Token     string `json:"token"`
	Timestamp string `json:"timestamp"`
}

// Models for Requests
type AccountRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}
type VerifyRequest struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

// Models for Responses
type AccountResponse struct {
	Result string
}

type LoginResponse struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Result string
}
