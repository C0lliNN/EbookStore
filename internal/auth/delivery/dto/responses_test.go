//go:build unit
// +build unit

package dto

import (
	"github.com/bxcodec/faker/v3"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestFromCredentials(t *testing.T) {
	token := faker.Jwt()

	credentials := model.Credentials{Token: token}

	expected := CredentialsResponse{Token: token}
	actual := FromCredentials(credentials)

	assert.Equal(t, expected, actual)
}
