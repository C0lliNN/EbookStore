package server

import (
	"time"

	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/gin-gonic/gin"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (e *LoggerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger := log.FromContext(c).With(
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", latency.String(),
		)

		logger.Info(path)
	}
}
