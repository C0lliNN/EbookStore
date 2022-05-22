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
	ServerTest
}

func TestAuthenticatorHandler(t *testing.T) {
	suite.Run(t, new(AuthenticatorHandlerTestSuite))
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
	require.NotNil(s.T(), err)
	require.Equal(s.T(), http.StatusBadRequest, response.StatusCode)

	errorResponse := &server.ErrorResponse{}
	err = json.NewDecoder(response.Body).Decode(errorResponse)

	require.Nil(s.T(), err)
	require.Equal(s.T(), "the payload is not valid", errorResponse.Message)
}

//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidLastName() {
//	password := faker.PASSWORD
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             "",
//		Email:                faker.Email(),
//		Password:             password,
//		PasswordConfirmation: password,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//
//	var user auth.User
//	err = s.db.First(&user).Error
//	assert.NotNil(s.T(), err)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidEmail() {
//	password := faker.PASSWORD
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             faker.LastName(),
//		Email:                "invalid-email",
//		Password:             password,
//		PasswordConfirmation: password,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//
//	var user auth.User
//	err = s.db.First(&user).Error
//	assert.NotNil(s.T(), err)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidPassword() {
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             faker.LastName(),
//		Email:                faker.Email(),
//		Password:             "1234",
//		PasswordConfirmation: "1234",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//
//	var user auth.User
//	err = s.db.First(&user).Error
//	assert.NotNil(s.T(), err)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_WithInvalidPasswordConfirmation() {
//	password := faker.PASSWORD
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             faker.LastName(),
//		Email:                faker.Email(),
//		Password:             password,
//		PasswordConfirmation: "12341234",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//
//	var user auth.User
//	err = s.db.First(&user).Error
//	assert.NotNil(s.T(), err)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_WithDuplicateEmail() {
//	persistedUser := factory.NewUser()
//	err := s.db.Create(persistedUser).Error
//	require.Nil(s.T(), err)
//
//	password := faker.PASSWORD
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             faker.LastName(),
//		Email:                persistedUser.Email,
//		Password:             password,
//		PasswordConfirmation: password,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//
//	var user auth.User
//	err = s.db.First(&user, "first_name = ? AND last_name = ?", payload.FirstName, payload.LastName).Error
//	assert.NotNil(s.T(), err)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestRegister_Successfully() {
//	password := faker.PASSWORD
//	payload := auth.RegisterRequest{
//		FirstName:            faker.FirstName(),
//		LastName:             faker.LastName(),
//		Email:                faker.Email(),
//		Password:             password,
//		PasswordConfirmation: password,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/register", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.register(s.context)
//
//	var credentials auth.CredentialsResponse
//	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
//	require.Nil(s.T(), err)
//
//	assert.Equal(s.T(), http.StatusCreated, s.recorder.Code)
//	assert.NotEmpty(s.T(), credentials.Token)
//
//	var user auth.User
//	err = s.db.First(&user).Error
//	require.Nil(s.T(), err)
//
//	assert.Len(s.T(), user.ID, 36)
//	assert.Equal(s.T(), payload.FirstName, user.FirstName)
//	assert.Equal(s.T(), payload.LastName, user.LastName)
//	assert.Equal(s.T(), payload.Email, user.Email)
//	assert.NotEmpty(s.T(), user.Password)
//	assert.Equal(s.T(), auth.Customer, user.Role)
//	assert.NotZero(s.T(), user.CreatedAt)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestLogin_WithMalFormedEmail() {
//	payload := auth.LoginRequest{
//		Email:    "invalid-email",
//		Password: "password",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.login(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestLogin_WithMalFormedPassword() {
//	payload := auth.LoginRequest{
//		Email:    faker.Email(),
//		Password: "1234",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.login(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestLogin_UnknownEmail() {
//	payload := auth.LoginRequest{
//		Email:    faker.Email(),
//		Password: "password",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.login(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestLogin_WithInvalidPassword() {
//	user := factory.NewUser()
//	err := s.db.Create(user).Error
//	require.Nil(s.T(), err)
//
//	payload := auth.LoginRequest{
//		Email:    user.Email,
//		Password: "wrong-password",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.login(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestLogin_Successfully() {
//	password := faker.PASSWORD
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
//	user := factory.NewUser()
//	user.Password = string(hashedPassword)
//
//	err := s.db.Create(user).Error
//	require.Nil(s.T(), err)
//
//	payload := auth.LoginRequest{
//		Email:    user.Email,
//		Password: password,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/login", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.login(s.context)
//
//	var credentials auth.CredentialsResponse
//	err = json.NewDecoder(s.recorder.Result().Body).Decode(&credentials)
//	require.Nil(s.T(), err)
//
//	assert.Equal(s.T(), http.StatusOK, s.recorder.Code)
//	assert.NotEmpty(s.T(), credentials.Token)
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestResetPassword_WithMalformedEmail() {
//	payload := auth.PasswordResetRequest{
//		Email: "invalid-email",
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.resetPassword(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestResetPassword_WithUnknownEmail() {
//	user := factory.NewUser()
//
//	payload := auth.PasswordResetRequest{
//		Email: user.Email,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.resetPassword(s.context)
//
//	assert.NotEmpty(s.T(), s.context.Errors.Errors())
//}
//
//func (s *AuthenticatorHandlerTestSuite) TestResetPassword_Successfully() {
//	user := factory.NewUser()
//	err := s.db.Create(user).Error
//
//	payload := auth.PasswordResetRequest{
//		Email: user.Email,
//	}
//
//	data, err := json.Marshal(payload)
//	require.Nil(s.T(), err)
//
//	request := httptest.NewRequest("POST", s.baseURL+"/password-reset", bytes.NewReader(data))
//	s.context.Request = request
//
//	s.handler.resetPassword(s.context)
//
//	assert.Empty(s.T(), s.context.Errors.Errors())
//
//	var updated auth.User
//	err = s.db.First(&updated, "id = ?", user.ID).Error
//	require.Nil(s.T(), err)
//
//	assert.NotEqual(s.T(), updated.Password, user.Password)
//}
