//go:build unit
// +build unit

package middleware

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/test/factory"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type AdminMiddlewareTestSuite struct {
	suite.Suite
	context    *gin.Context
	middleware AdminMiddleware
}

func (s *AdminMiddlewareTestSuite) SetupTest() {
	s.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	s.middleware = AdminMiddleware{}
}

func TestAdminMiddlewareRun(t *testing.T) {
	suite.Run(t, new(AdminMiddlewareTestSuite))
}

func (s *AdminMiddlewareTestSuite) TestHandler_WhenUserIsNotSet() {
	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())
}

func (s *AdminMiddlewareTestSuite) TestHandler_WhenUserIsCustomer() {
	user := factory.NewUser()
	user.Role = model.Customer

	s.context.Set("user", user)
	s.middleware.Handler()(s.context)

	assert.True(s.T(), s.context.IsAborted())
}

func (s *AdminMiddlewareTestSuite) TestHandler_WhenUserIsAdmin() {
	user := factory.NewUser()
	user.Role = model.Admin

	s.context.Set("user", user)
	s.middleware.Handler()(s.context)

	assert.False(s.T(), s.context.IsAborted())
}
