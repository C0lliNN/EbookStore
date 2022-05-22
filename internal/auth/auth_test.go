package auth_test

import (
	"context"
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
	mocks2 "github.com/c0llinn/ebook-store/internal/mocks/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	newIdMethod                  = "NewID"
	saveMethod                   = "Save"
	findByEmail                  = "FindByEmail"
	updateMethod                 = "Update"
	generateTokenMethod          = "GenerateTokenForUser"
	newPasswordMethod            = "NewPassword"
	sendEmailMethod              = "SendPasswordResetEmail"
	hashPasswordMethod           = "HashPassword"
	compareHashAndPasswordMethod = "CompareHashAndPassword"
	validateMethod               = "Validate"
)

type AuthenticatorTestSuite struct {
	suite.Suite
	token             *mocks2.TokenHandler
	repo              *mocks2.Repository
	emailClient       *mocks2.EmailClient
	passwordGenerator *mocks2.PasswordGenerator
	hash              *mocks2.HashHandler
	idGenerator       *mocks2.IDGenerator
	validator         *mocks2.Validator
	authenticator     *auth.Authenticator
}

func (s *AuthenticatorTestSuite) SetupTest() {
	s.token = new(mocks2.TokenHandler)
	s.repo = new(mocks2.Repository)
	s.emailClient = new(mocks2.EmailClient)
	s.passwordGenerator = new(mocks2.PasswordGenerator)
	s.hash = new(mocks2.HashHandler)
	s.idGenerator = new(mocks2.IDGenerator)
	s.validator = new(mocks2.Validator)

	config := auth.Config{
		Repository:        s.repo,
		Tokener:           s.token,
		Hasher:            s.hash,
		EmailClient:       s.emailClient,
		PasswordGenerator: s.passwordGenerator,
		IDGenerator:       s.idGenerator,
		Validator:         s.validator,
	}

	s.authenticator = auth.New(config)
}

func TestAuthenticator(t *testing.T) {
	suite.Run(t, new(AuthenticatorTestSuite))
}

func (s *AuthenticatorTestSuite) TestRegister_WhenPasswordValidationFails() {
	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}

	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 0)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 0)
	s.repo.AssertNotCalled(s.T(), saveMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestRegister_WhenPasswordHashingFails() {
	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	s.validator.On(validateMethod, request).Return(nil)

	s.idGenerator.On(newIdMethod).Return("user-id")
	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("", fmt.Errorf("some-error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some-error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), saveMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestRegister_WhenRepositoryFails() {
	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	s.validator.On(validateMethod, request).Return(nil)

	s.idGenerator.On(newIdMethod).Return("user-id")
	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(fmt.Errorf("some error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestRegister_WhenTokenGenerationFails() {
	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	s.validator.On(validateMethod, request).Return(nil)

	s.idGenerator.On(newIdMethod).Return("user-id")
	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(nil)

	s.token.On(generateTokenMethod, updatedUser).Return("", fmt.Errorf("some error"))

	_, err := s.authenticator.Register(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthenticatorTestSuite) TestRegister_Successfully() {
	request := auth.RegisterRequest{
		FirstName:            "Raphael",
		LastName:             "Collin",
		Email:                "raphael@test.com",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	s.validator.On(validateMethod, request).Return(nil)

	s.idGenerator.On(newIdMethod).Return("user-id")
	user := request.User("user-id")
	s.hash.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, context.TODO(), &updatedUser).Return(nil)
	s.token.On(generateTokenMethod, updatedUser).Return("token", nil)

	response, err := s.authenticator.Register(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), auth.NewCredentialsResponse(auth.Credentials{Token: "token"}), response)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.idGenerator.AssertNumberOfCalls(s.T(), newIdMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthenticatorTestSuite) TestLogin_WhenValidationFails() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}
	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	_, err := s.authenticator.Login(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 0)
	s.hash.AssertNotCalled(s.T(), compareHashAndPasswordMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestLogin_WhenUserWasNotFound() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}
	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, fmt.Errorf("some error"))

	_, err := s.authenticator.Login(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNotCalled(s.T(), compareHashAndPasswordMethod)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestLogin_WhenPasswordsDontMatch() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}

	user := auth.User{ID: "some-id", Email: request.Email, Password: "some-password"}
	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.hash.On(compareHashAndPasswordMethod, user.Password, request.Password).Return(auth.ErrWrongPassword)

	_, err := s.authenticator.Login(context.TODO(), request)

	assert.Equal(s.T(), auth.ErrWrongPassword, err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.token.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthenticatorTestSuite) TestLogin_Successfully() {
	request := auth.LoginRequest{
		Email:    "email@test.com",
		Password: "12345678",
	}

	user := auth.User{ID: "some-id", Email: request.Email, Password: "some-password"}

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.hash.On(compareHashAndPasswordMethod, user.Password, request.Password).Return(nil)
	s.token.On(generateTokenMethod, user).Return("token", nil)

	response, err := s.authenticator.Login(context.TODO(), request)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), auth.NewCredentialsResponse(auth.Credentials{Token: "token"}), response)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.token.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenValidationFails() {
	request := auth.PasswordResetRequest{Email: "some email"}
	s.validator.On(validateMethod, request).Return(fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 0)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.hash.AssertNotCalled(s.T(), hashPasswordMethod)
	s.passwordGenerator.AssertNotCalled(s.T(), newPasswordMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenUserWasNotFound() {
	request := auth.PasswordResetRequest{Email: "some email"}
	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.hash.AssertNotCalled(s.T(), hashPasswordMethod)
	s.passwordGenerator.AssertNotCalled(s.T(), newPasswordMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenPasswordHashingFails() {
	request := auth.PasswordResetRequest{Email: "some email"}

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), request.Email).Return(auth.User{}, nil)
	s.passwordGenerator.On(newPasswordMethod).Return("new-password")
	s.hash.On(hashPasswordMethod, "new-password").Return("", fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenUpdateFails() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthenticatorTestSuite) TestResetPassword_WhenEmailSendingFails() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(nil)
	s.emailClient.On(sendEmailMethod, context.TODO(), user, newPassword).Return(fmt.Errorf("some error"))

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}

func (s AuthenticatorTestSuite) TestResetPassword_Successfully() {
	request := auth.PasswordResetRequest{Email: "some email"}
	user := auth.User{Email: request.Email, Password: "another-password"}
	newPassword := "password"

	s.validator.On(validateMethod, request).Return(nil)
	s.repo.On(findByEmail, context.TODO(), user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.hash.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, context.TODO(), &user).Return(nil)
	s.emailClient.On(sendEmailMethod, context.TODO(), user, newPassword).Return(nil)

	err := s.authenticator.ResetPassword(context.TODO(), request)

	assert.Nil(s.T(), err)

	s.validator.AssertNumberOfCalls(s.T(), validateMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.hash.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}
