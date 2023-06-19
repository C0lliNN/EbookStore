package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"time"
)

// RateLimitMiddleware a simple wrapper around the ulele/limiter. "Wrap the code you control around the code you don't control"
type RateLimitMiddleware struct {
	rate  limiter.Rate
	store limiter.Store
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rate: limiter.Rate{
			Period: 1 * time.Hour,
			Limit:  1000,
		},
		store: memory.NewStore(),
	}
}

func (m *RateLimitMiddleware) Handler() gin.HandlerFunc {
	return mgin.NewMiddleware(limiter.New(m.store, m.rate))
}
