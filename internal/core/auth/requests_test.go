package auth

import (
	"github.com/bxcodec/faker/v3"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRegisterRequest_User(t *testing.T) {
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

	expected := User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Role:      Customer,
		Email:     email,
		Password:  password,
	}

	actual := registerRequest.User(id)

	assert.Equal(t, expected, actual)
}
