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
	generateTokenMethod = "GenerateTokenForUser"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	jwt     *mock.JWTWrapper
	repo    *mock.UserRepository
	useCase AuthUseCase
}

func (s *AuthUseCaseTestSuite) SetupTest() {
	s.jwt = new(mock.JWTWrapper)
	s.repo = new(mock.UserRepository)
	s.useCase = AuthUseCase{jwt: s.jwt, repo: s.repo}
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
