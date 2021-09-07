// +build unit

package dto

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRegisterRequest_ToDomain(t *testing.T) {
	id := faker.UUIDHyphenated()
	firstName := faker.FirstName()
	lastName := faker.LastName()
	email := faker.Email()
	password := faker.Password()
	confirmPassword := password

	registerRequest := RegisterRequest{
		FirstName:            firstName,
		LastName:             lastName,
		Email:                email,
		Password:             password,
		PasswordConfirmation: confirmPassword,
	}

	expected := model.User{
		ID: id,
		FirstName: firstName,
		LastName: lastName,
		Role: model.Customer,
		Email: email,
		Password: password,
	}

	actual := registerRequest.ToDomain(id)

	assert.Equal(t, expected, actual)
}
