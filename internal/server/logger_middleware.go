package server

import (
	"github.com/c0llinn/ebook-store/internal/log"
	"github.com/gin-gonic/gin"
	"time"
)

type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (e *LoggerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.FromContext(c)

		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			logger := logger.With(
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
}
