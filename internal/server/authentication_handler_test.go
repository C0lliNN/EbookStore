package server_test

import (
	"bytes"
	"encoding/json"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/server"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthenticatorHandlerTestSuite struct {
	ServerSuiteTest
}

func TestAuthenticatorHandler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(AuthenticatorHandlerTestSuite))
}

func (s *AuthenticatorHandlerTestSuite) TearDownTest() {
	s.db.Delete(&auth.User{}, "1 = 1")
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_Successfully() {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, response.StatusCode)

	credentials := &auth.CredentialsResponse{}
	err = json.NewDecoder(response.Body).Decode(credentials)

	require.Nil(s.T(), err)
	require.NotEmpty(s.T(), credentials.Token)
}

func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidData() {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusBadRequest, response.StatusCode)

	errorResponse := &server.ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(errorResponse)

	require.Nil(s.T(), err)
	require.Equal(s.T(), "the payload is not valid", errorResponse.Message)
}

func (s *AuthenticatorHandlerTestSuite) TestLogin_Failure() {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, response.StatusCode)

	payload2 := auth.LoginRequest{
		Email:    "raphael@test.com",
		Password: "password",
	}

	data, err = json.Marshal(payload2)
	require.Nil(s.T(), err)

	response, err = http.Post(s.baseURL+"/login", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)

	credentials := &auth.CredentialsResponse{}
	err = json.NewDecoder(response.Body).Decode(credentials)

	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusOK, response.StatusCode)
}

func (s *AuthenticatorHandlerTestSuite) TestLogin_Success() {
	s.createUser()

	payload := auth.LoginRequest{
		Email:    "raphael@test.com",
		Password: "password",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/login", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)

	credentials := &auth.CredentialsResponse{}
	err = json.NewDecoder(response.Body).Decode(credentials)
	require.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusOK, response.StatusCode)
	require.NotEmpty(s.T(), credentials.Token)
}

func (s *AuthenticatorHandlerTestSuite) TestResetPassword_Failure() {
	payload := auth.PasswordResetRequest{
		Email:    "raphael@test.com",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/password-reset", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)

	errorResponse := &server.ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(errorResponse)
	require.Nil(s.T(), err)

	require.Equal(s.T(), http.StatusNotFound, response.StatusCode)
	require.Equal(s.T(), "the provided User was not found", errorResponse.Message)
}

func (s *AuthenticatorHandlerTestSuite) TestResetPassword_Success() {
	s.createUser()
	payload := auth.PasswordResetRequest{
		Email:    "raphael@test.com",
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/password-reset", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusNoContent, response.StatusCode)
}

func (s *AuthenticatorHandlerTestSuite) createUser() {
	password := "password"
	payload := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             password,
		PasswordConfirmation: password,
	}

	data, err := json.Marshal(payload)
	require.Nil(s.T(), err)

	response, err := http.Post(s.baseURL+"/register", "application/json", bytes.NewReader(data))
	require.Nil(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, response.StatusCode)
}
