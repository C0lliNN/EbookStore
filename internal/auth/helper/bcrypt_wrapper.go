package helper

import "golang.org/x/crypto/bcrypt"

const (
	cost = 12
)

type BcryptWrapper struct{}

func NewBcryptWrapper() BcryptWrapper {
	return BcryptWrapper{}
}

func (w BcryptWrapper) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func (w BcryptWrapper) CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
