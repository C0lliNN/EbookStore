package auth

type CredentialsResponse struct {
	Token string `json:"token"`
}

func NewCredentialsResponse(credentials Credentials) CredentialsResponse {
	return CredentialsResponse(credentials)
}
