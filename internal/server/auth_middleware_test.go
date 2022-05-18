package server_test

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/server"
	mocks "github.com/c0llinn/ebook-store/mocks/server"

	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"strings"
	"testing"
)

const extractUserFromTokenMethod = "ExtractUserFromToken"

type AuthMiddlewareTestSuite struct {
	suite.Suite
	context    *gin.Context
	token      *mocks.TokenHandler
	middleware *server.AuthenticationMiddleware
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.context.Request = httptest.NewRequest("GET", "/books", strings.NewReader(""))

	s.token = new(mocks.TokenHandler)
	s.middleware = server.NewAuthenticationMiddleware(s.token)
}

func TestAuthMiddlewareRun(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithoutHeader() {
	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.token.AssertNotCalled(s.T(), extractUserFromTokenMethod)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithMalformedHeader() {
	s.context.Request.Header.Set("Authorization", "token")

	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.token.AssertNotCalled(s.T(), extractUserFromTokenMethod)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithInvalidToken() {
	s.context.Request.Header.Set("Authorization", "Bearer token")

	s.token.On(extractUserFromTokenMethod, "token").Return(auth.User{}, fmt.Errorf("some error"))

	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.token.AssertNumberOfCalls(s.T(), extractUserFromTokenMethod, 1)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithValidToken() {
	s.context.Request.Header.Set("Authorization", "Bearer token")

	user := factory.NewUser()
	s.token.On(extractUserFromTokenMethod, "token").Return(user, nil)

	s.middleware.Handler()(s.context)
	assert.False(s.T(), s.context.IsAborted())

	actual, _ := s.context.Get("user")
	assert.Equal(s.T(), user, actual.(auth.User))

	s.token.AssertNumberOfCalls(s.T(), extractUserFromTokenMethod, 1)
}
