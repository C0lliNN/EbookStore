// +build unit

package token

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JWTWrapperTestSuite struct {
	suite.Suite
	secret []byte
	jwtWrapper JWTWrapper
}

func (s *JWTWrapperTestSuite) SetupTest() {
	s.secret = []byte("secret")
	s.jwtWrapper = JWTWrapper{secret: s.secret}
}

func TestJWTWrapperRun(t *testing.T) {
	suite.Run(t, new(JWTWrapperTestSuite))
}

func (s *JWTWrapperTestSuite) TestGenerateTokenForUser() {
	user := model.User{
		ID: "some-id",
		Email: "test@test.com",
		FirstName: "first",
		LastName: "last",
		Role: model.Admin,
	}

	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiaWQiOiJzb21lLWlkIiwibmFtZSI6ImZpcnN0IGxhc3QifQ.OLYJtZCJzlKbzJ9jXRrY9cjndGMItSrYIWij2bFnevI"
	actual, err := s.jwtWrapper.GenerateTokenForUser(user)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
}

func (s *JWTWrapperTestSuite) TestExtractUserFromToken() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiaWQiOiJzb21lLWlkIiwibmFtZSI6ImZpcnN0IGxhc3QifQ.OLYJtZCJzlKbzJ9jXRrY9cjndGMItSrYIWij2bFnevI"

	expected := model.User{
		ID: "some-id",
		Email: "test@test.com",
		FirstName: "first",
		LastName: "last",
		Role: model.Admin,
	}
	actual, err := s.jwtWrapper.ExtractUserFromToken(token)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expected, actual)
}