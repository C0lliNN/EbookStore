package hash

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcryptWrapper_HashPassword(t *testing.T) {
	hashedPassword, err := NewBcryptWrapper().HashPassword("some-password")

	assert.Nil(t, err)
	assert.Len(t, hashedPassword, 60)
}

func TestBcryptWrapper_CompareHashAndPasswordWithInvalidPassword(t *testing.T) {
	hashedPassword := "some-hashed-password"
	password := "password"

	err := NewBcryptWrapper().CompareHashAndPassword(hashedPassword, password)

	assert.NotNil(t, err)
}

func TestBcryptWrapper_CompareHashAndPasswordWithValidPassword(t *testing.T) {
	password := "some-password"

	hashedPassword, err := NewBcryptWrapper().HashPassword(password)
	require.Nil(t, err)

	err = NewBcryptWrapper().CompareHashAndPassword(hashedPassword, password)
	assert.Nil(t, err)
}
