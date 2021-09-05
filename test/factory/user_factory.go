package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewUser() auth.User {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), 12)

	return auth.User{
		ID:        faker.UUIDHyphenated(),
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Role:      auth.Customer,
		Password:  string(bytes),
		CreatedAt: time.Now().Unix(),
	}
}
