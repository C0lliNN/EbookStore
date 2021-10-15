//go:build unit
// +build unit

package usecase

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth/mock"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/common"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

const (
	saveMethod                   = "Save"
	findByEmail                  = "FindByEmail"
	updateMethod                 = "Update"
	generateTokenMethod          = "GenerateTokenForUser"
	newPasswordMethod            = "NewPassword"
	sendEmailMethod              = "SendPasswordResetEmail"
	hashPasswordMethod           = "HashPassword"
	compareHashAndPasswordMethod = "CompareHashAndPassword"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	jwt               *mock.JWTWrapper
	repo              *mock.UserRepository
	emailClient       *mock.EmailClient
	passwordGenerator *mock.PasswordGenerator
	bcrypt            *mock.BcryptWrapper
	useCase           AuthUseCase
}

func (s *AuthUseCaseTestSuite) SetupTest() {
	s.jwt = new(mock.JWTWrapper)
	s.repo = new(mock.UserRepository)
	s.emailClient = new(mock.EmailClient)
	s.passwordGenerator = new(mock.PasswordGenerator)
	s.bcrypt = new(mock.BcryptWrapper)

	s.useCase = AuthUseCase{jwt: s.jwt, repo: s.repo, emailClient: s.emailClient, passwordGenerator: s.passwordGenerator, bcrypt: s.bcrypt}
}

func TestAuthUseCaseRun(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenPasswordHashingFails() {
	user := factory.NewUser()
	s.bcrypt.On(hashPasswordMethod, user.Password).Return("", fmt.Errorf("some-error"))

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), fmt.Errorf("some-error"), err)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), saveMethod)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenRepositoryFails() {
	user := factory.NewUser()
	s.bcrypt.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, &updatedUser).Return(gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenTokenGenerationFails() {
	user := factory.NewUser()
	s.bcrypt.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, &updatedUser).Return(nil)

	s.jwt.On(generateTokenMethod, updatedUser).Return("", gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenNoErrorWasFound() {
	user := factory.NewUser()
	s.bcrypt.On(hashPasswordMethod, user.Password).Return("hashed-password", nil)

	updatedUser := user
	updatedUser.Password = "hashed-password"
	s.repo.On(saveMethod, &updatedUser).Return(nil)
	s.jwt.On(generateTokenMethod, updatedUser).Return("token", nil)

	credentials, err := s.useCase.Register(user)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.Credentials{Token: "token"}, credentials)

	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthUseCaseTestSuite) TestLogin_WhenUserWasNotFound() {
	s.repo.On(findByEmail, "email@test.com").Return(model.User{}, fmt.Errorf("some error"))

	_, err := s.useCase.Login("email@test.com", "password")

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNotCalled(s.T(), compareHashAndPasswordMethod)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestLogin_WhenPasswordsDontMatch() {
	user := factory.NewUser()

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.bcrypt.On(compareHashAndPasswordMethod, user.Password, "password").Return(&common.ErrWrongPassword{})

	_, err := s.useCase.Login(user.Email, "password")

	assert.IsType(s.T(), &common.ErrWrongPassword{}, err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestLogin_Successfully() {
	user := factory.NewUser()

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.bcrypt.On(compareHashAndPasswordMethod, user.Password, "password").Return(nil)
	s.jwt.On(generateTokenMethod, user).Return("token", nil)

	credentials, err := s.useCase.Login(user.Email, "password")

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.Credentials{Token: "token"}, credentials)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), compareHashAndPasswordMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenUserWasNotFound() {
	email := faker.Email()
	s.repo.On(findByEmail, email).Return(model.User{}, fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.bcrypt.AssertNotCalled(s.T(), hashPasswordMethod)
	s.passwordGenerator.AssertNotCalled(s.T(), newPasswordMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenPasswordHashingFails() {
	email := faker.Email()
	s.repo.On(findByEmail, email).Return(model.User{}, nil)
	s.passwordGenerator.On(newPasswordMethod).Return("new-password")
	s.bcrypt.On(hashPasswordMethod, "new-password").Return("", fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenUpdateFails() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.bcrypt.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, &user).Return(fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(user.Email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenEmailSendingFails() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.bcrypt.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, &user).Return(nil)
	s.emailClient.On(sendEmailMethod, user, newPassword).Return(fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(user.Email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}

func (s AuthUseCaseTestSuite) TestResetPassword_Successfully() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.bcrypt.On(hashPasswordMethod, newPassword).Return("new-hashed-password", nil)
	user.Password = "new-hashed-password"

	s.repo.On(updateMethod, &user).Return(nil)
	s.emailClient.On(sendEmailMethod, user, newPassword).Return(nil)

	err := s.useCase.ResetPassword(user.Email)

	assert.Nil(s.T(), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.bcrypt.AssertNumberOfCalls(s.T(), hashPasswordMethod, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}
