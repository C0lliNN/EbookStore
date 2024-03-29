package server_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/platform/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const extractUserFromTokenMethod = "ExtractUserFromToken"

type AuthMiddlewareTestSuite struct {
	suite.Suite
	context    *gin.Context
	token      *server.MockTokenHandler
	middleware *server.AuthenticationMiddleware
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.context.Request = httptest.NewRequest("GET", "/books", strings.NewReader(""))

	s.token = new(server.MockTokenHandler)
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

	user := auth.User{
		ID:        "some-id",
		FirstName: "Raphael",
		LastName:  "Collin",
		Email:     "raphael@test.com",
		Role:      auth.Customer,
		Password:  "password",
		CreatedAt: time.Now().Unix(),
	}
	s.token.On(extractUserFromTokenMethod, "token").Return(user, nil)

	s.middleware.Handler()(s.context)
	assert.False(s.T(), s.context.IsAborted())

	assert.Equal(s.T(), user.ID, s.context.Value("userId").(string))
	assert.Equal(s.T(), user.IsAdmin(), s.context.Value("admin").(bool))

	s.token.AssertNumberOfCalls(s.T(), extractUserFromTokenMethod, 1)
}
