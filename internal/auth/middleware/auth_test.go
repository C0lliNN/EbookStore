package middleware

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth/mock"
	"github.com/c0llinn/ebook-store/internal/auth/model"
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
	context *gin.Context
	jwt *mock.JWTWrapper
	middleware AuthenticationMiddleware
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.context.Request = httptest.NewRequest("GET", "/books", strings.NewReader(""))

	s.jwt = new(mock.JWTWrapper)
	s.middleware = NewAuthenticationMiddleware(s.jwt)
}

func TestAuthMiddlewareRun(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithoutHeader() {
	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.jwt.AssertNotCalled(s.T(), extractUserFromTokenMethod)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithMalformedHeader() {
	s.context.Request.Header.Set("Authorization", "token")

	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.jwt.AssertNotCalled(s.T(), extractUserFromTokenMethod)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithInvalidToken() {
	s.context.Request.Header.Set("Authorization", "Bearer token")

	s.jwt.On(extractUserFromTokenMethod, "token").Return(model.User{}, fmt.Errorf("some error"))

	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())

	s.jwt.AssertNumberOfCalls(s.T(), extractUserFromTokenMethod, 1)
}

func (s *AuthMiddlewareTestSuite) TestHandler_WithValidToken() {
	s.context.Request.Header.Set("Authorization", "Bearer token")

	user := factory.NewUser()
	s.jwt.On(extractUserFromTokenMethod, "token").Return(user, nil)

	s.middleware.Handler()(s.context)
	assert.False(s.T(), s.context.IsAborted())

	actual, _ := s.context.Get("user")
	assert.Equal(s.T(), user, actual.(model.User))

	s.jwt.AssertNumberOfCalls(s.T(), extractUserFromTokenMethod, 1)
}