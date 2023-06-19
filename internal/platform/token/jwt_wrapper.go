package token

import (
	"fmt"
	"strings"

	"github.com/ebookstore/internal/core/auth"
	"github.com/golang-jwt/jwt"
)

type HMACSecret []byte

type JWTWrapper struct {
	secret HMACSecret
}

func NewJWTWrapper(secret HMACSecret) *JWTWrapper {
	return &JWTWrapper{secret: secret}
}

func (w *JWTWrapper) GenerateTokenForUser(user auth.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.FullName(),
		"admin": user.IsAdmin(),
	})

	signedString, err := token.SignedString([]byte(w.secret))
	if err != nil {
		return "", fmt.Errorf("(GenerateTokenForUser) failed generating token for user: %w", err)
	}

	return signedString, nil
}

func (w *JWTWrapper) ExtractUserFromToken(tokenString string) (auth.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("(ExtractUserFromToken) unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(w.secret), nil
	})

	if err != nil {
		return auth.User{}, fmt.Errorf("(ExtractUserFromToken) failed parsing jwt token")
	}

	user := auth.User{}
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

	return user, nil
}
