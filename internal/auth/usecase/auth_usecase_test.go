// +build unit

package usecase

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth/mock"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

const (
	saveMethod          = "Save"
	findByEmail         = "FindByEmail"
	updateMethod = "Update"
	generateTokenMethod = "GenerateTokenForUser"
	newPasswordMethod = "NewPassword"
	sendEmailMethod = "SendPasswordResetEmail"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	jwt               *mock.JWTWrapper
	repo              *mock.UserRepository
	emailClient       *mock.EmailClient
	passwordGenerator *mock.PasswordGenerator
	useCase           AuthUseCase
}

func (s *AuthUseCaseTestSuite) SetupTest() {
	s.jwt = new(mock.JWTWrapper)
	s.repo = new(mock.UserRepository)
	s.emailClient = new(mock.EmailClient)
	s.passwordGenerator = new(mock.PasswordGenerator)

	s.useCase = AuthUseCase{jwt: s.jwt, repo: s.repo, emailClient: s.emailClient, passwordGenerator: s.passwordGenerator}
}

func TestAuthUseCaseRun(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenRepositoryFails() {
	user := factory.NewUser()

	s.repo.On(saveMethod, mock2.AnythingOfType("*model.User")).Return(gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenTokenGenerationFails() {
	user := factory.NewUser()

	s.repo.On(saveMethod, mock2.AnythingOfType("*model.User")).Return(nil)
	s.jwt.On(generateTokenMethod, mock2.AnythingOfType("model.User")).Return("", gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenNoErrorWasFound() {
	user := factory.NewUser()

	s.repo.On(saveMethod, mock2.AnythingOfType("*model.User")).Return(nil)
	s.jwt.On(generateTokenMethod, mock2.AnythingOfType("model.User")).Return("token", nil)

	credentials, err := s.useCase.Register(user)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.Credentials{Token: "token"}, credentials)

	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthUseCaseTestSuite) TestLogin_WhenUserWasNotFound() {
	s.repo.On(findByEmail, "email@test.com").Return(model.User{}, fmt.Errorf("some error"))

	_, err := s.useCase.Login("email@test.com", "password")

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestLogin_WhenPasswordsDontMatch() {
	user := factory.NewUser()

	s.repo.On(findByEmail, user.Email).Return(user, nil)

	_, err := s.useCase.Login(user.Email, "password")

	assert.IsType(s.T(), &model.ErrWrongPassword{}, err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestLogin_Successfully() {
	password := faker.Password()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user := factory.NewUser()
	user.Password = string(hashedPassword)

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.jwt.On(generateTokenMethod, user).Return("token", nil)

	credentials, err := s.useCase.Login(user.Email, password)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.Credentials{Token: "token"}, credentials)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenUserWasNotFound() {
	email := faker.Email()
	s.repo.On(findByEmail, email).Return(model.User{}, fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNotCalled(s.T(), updateMethod)
	s.passwordGenerator.AssertNotCalled(s.T(), newPasswordMethod)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenUpdateFails() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.repo.On(updateMethod, mock2.Anything).Return(fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(user.Email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNotCalled(s.T(), sendEmailMethod)
}

func (s AuthUseCaseTestSuite) TestResetPassword_WhenEmailSendingFails() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.repo.On(updateMethod, mock2.Anything).Return(nil)
	s.emailClient.On(sendEmailMethod, mock2.Anything, newPassword).Return(fmt.Errorf("some error"))

	err := s.useCase.ResetPassword(user.Email)

	assert.Equal(s.T(), fmt.Errorf("some error"), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}

func (s AuthUseCaseTestSuite) TestResetPassword_Successfully() {
	user := factory.NewUser()
	newPassword := "password"

	s.repo.On(findByEmail, user.Email).Return(user, nil)
	s.passwordGenerator.On(newPasswordMethod).Return(newPassword)
	s.repo.On(updateMethod, mock2.Anything).Return(nil)
	s.emailClient.On(sendEmailMethod, mock2.Anything, newPassword).Return(nil)

	err := s.useCase.ResetPassword(user.Email)

	assert.Nil(s.T(), err)

	s.repo.AssertNumberOfCalls(s.T(), findByEmail, 1)
	s.repo.AssertNumberOfCalls(s.T(), updateMethod, 1)
	s.passwordGenerator.AssertNumberOfCalls(s.T(), newPasswordMethod, 1)
	s.emailClient.AssertNumberOfCalls(s.T(), sendEmailMethod, 1)
}
