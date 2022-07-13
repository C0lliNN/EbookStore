package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimitMiddleware(t *testing.T) {
	m := NewRateLimitMiddleware()

	assert.NotPanics(t, func() {
		h := m.Handler()
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("GET", "google.com", bytes.NewReader([]byte("")))
		h(ctx)
	})
}
