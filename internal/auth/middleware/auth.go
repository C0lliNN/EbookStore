package middleware

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type JWTWrapper interface {
	ExtractUserFromToken(token string) (model.User, error)
}

type AuthenticationMiddleware struct {
	jwt JWTWrapper
}

func NewAuthenticationMiddleware(jwt JWTWrapper) AuthenticationMiddleware {
	return AuthenticationMiddleware{jwt}
}

func (m AuthenticationMiddleware) Handler() gin.HandlerFunc {
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
		user, err := m.jwt.ExtractUserFromToken(token)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "You are not authorized.",
				"details": "The Bearer token is not valid",
			})
			return
		}

		context.Set("user", user)
		context.Next()
	}
}
