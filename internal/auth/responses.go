package auth


type CredentialsResponse struct {
	Token string `json:"Tokener"`
}

func FromCredentials(credentials Credentials) CredentialsResponse {
	return CredentialsResponse{Token: credentials.Token}
}
