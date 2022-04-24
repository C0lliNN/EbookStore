package auth

import (
	"github.com/bxcodec/faker/v3"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFromCredentials(t *testing.T) {
	token := faker.Jwt()

	credentials := Credentials{Token: token}

	expected := CredentialsResponse{Token: token}
	actual := FromCredentials(credentials)

	assert.Equal(t, expected, actual)
}
