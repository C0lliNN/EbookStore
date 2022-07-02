package server

import (
	"net/http"
	"regexp"

	"github.com/ebookstore/internal/auth"
	"github.com/gin-gonic/gin"
)

type TokenHandler interface {
	ExtractUserFromToken(token string) (auth.User, error)
}

type AuthenticationMiddleware struct {
	token TokenHandler
}

func NewAuthenticationMiddleware(token TokenHandler) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{token}
}

func (m *AuthenticationMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		header := context.GetHeader("Authorization")
		if header == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized.",
				"details": "The 'Authorization' header must be provided",
			})
			return
		}

		if match, err := regexp.MatchString("Bearer .+", header); !match || err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized.",
				"details": "The 'Authorization' header must in be in the format 'Bearer token'",
			})
			return
		}

		token := header[7:]
		user, err := m.token.ExtractUserFromToken(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized.",
				"details": "The Bearer token is not valid",
			})
			return
		}

		context.Set("userId", user.ID)
		context.Set("admin", user.IsAdmin())
		context.Next()
	}
}
