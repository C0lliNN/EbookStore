// +build unit

package usecase

import (
	"github.com/c0llinn/ebook-store/internal/auth/mock"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

const (
	saveMethod = "Save"
	generateTokenMethod = "GenerateTokenForUser"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	jwt *mock.JWTWrapper
	repo *mock.UserRepository
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
	s.repo.On(saveMethod, &user).Return(gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNotCalled(s.T(), generateTokenMethod)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenTokenGenerationFails() {
	user := factory.NewUser()
	s.repo.On(saveMethod, &user).Return(nil)
	s.jwt.On(generateTokenMethod, user).Return("", gorm.ErrInvalidValue)

	_, err := s.useCase.Register(user)

	assert.Equal(s.T(), gorm.ErrInvalidValue, err)
	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}

func (s *AuthUseCaseTestSuite) TestRegister_WhenNoErrorWasFound() {
	user := factory.NewUser()
	s.repo.On(saveMethod, &user).Return(nil)
	s.jwt.On(generateTokenMethod, user).Return("token", nil)

	credentials, err := s.useCase.Register(user)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), model.Credentials{Token: "token"}, credentials)

	s.repo.AssertNumberOfCalls(s.T(), saveMethod, 1)
	s.jwt.AssertNumberOfCalls(s.T(), generateTokenMethod, 1)
}
