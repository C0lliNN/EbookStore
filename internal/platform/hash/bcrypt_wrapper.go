package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	cost = 12
)

type BcryptWrapper struct{}

func NewBcryptWrapper() *BcryptWrapper {
	return &BcryptWrapper{}
}

func (w *BcryptWrapper) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("(HashPassword) failed generating hash for the password: %w", err)
	}
	return string(bytes), nil
}

func (w *BcryptWrapper) CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("(CompareHashAndPassword) failed comparing hash with password: %w", err)
	}
	return nil
}
