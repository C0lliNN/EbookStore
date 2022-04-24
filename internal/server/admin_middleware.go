package server

import (
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() AdminMiddleware {
	return AdminMiddleware{}
}

func (m AdminMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		value, exists := context.Get("user")
		if !exists {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized.",
				"details": "The user must be authenticated",
			})
			return
		}

		user := value.(auth.User)
		if !user.IsAdmin() {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You cannot access this resource.",
				"details": "This resource is reserved for administrators.",
			})
			return
		}

		context.Next()
	}
}
