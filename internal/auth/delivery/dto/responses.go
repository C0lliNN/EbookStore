package dto

import "github.com/c0llinn/ebook-store/internal/auth/model"

type CredentialsResponse struct {
	Token string `json:"token"`
}

func FromCredentials(credentials model.Credentials) CredentialsResponse {
	return CredentialsResponse{Token: credentials.Token}
}
