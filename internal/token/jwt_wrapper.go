package token

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"strings"
)

type HMACSecret []byte

func NewHMACSecret() HMACSecret {
	return []byte(viper.GetString("JWT_SECRET"))
}

type JWTWrapper struct {
	secret HMACSecret
}

func NewJWTWrapper(secret HMACSecret) JWTWrapper {
	return JWTWrapper{secret: secret}
}

func (w JWTWrapper) GenerateTokenForUser(user auth.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.FullName(),
		"admin": user.IsAdmin(),
	})

	return token.SignedString([]byte(w.secret))
}

func (w JWTWrapper) ExtractUserFromToken(tokenString string) (user auth.User, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(w.secret), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.ID = claims["id"].(string)
		user.FirstName = strings.Split(claims["name"].(string), " ")[0]
		user.LastName = strings.Split(claims["name"].(string), " ")[1]
		user.Email = claims["email"].(string)

		if claims["admin"].(bool) {
			user.Role = auth.Admin
		} else {
			user.Role = auth.Customer
		}
	}

	return
}
