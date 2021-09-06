package factory

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewUser() model.User {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(faker.Password()), 12)

	return model.User{
		ID:        faker.UUIDHyphenated(),
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Email:     faker.Email(),
		Role:      model.Customer,
		Password:  string(bytes),
		CreatedAt: time.Now().Unix(),
	}
}
