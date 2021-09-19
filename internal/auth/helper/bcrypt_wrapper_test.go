// +build unit

package helper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcryptWrapper_HashPassword(t *testing.T) {
	hashedPassword, err := BcryptWrapper{}.HashPassword("some-password")

	assert.Nil(t, err)
	assert.Len(t, hashedPassword, 60)
}

func TestBcryptWrapper_CompareHashAndPasswordWithInvalidPassword(t *testing.T) {
	hashedPassword := "some-hashed-password"
	password := "password"

	err := BcryptWrapper{}.CompareHashAndPassword(hashedPassword, password)

	assert.NotNil(t, err)
}

func TestBcryptWrapper_CompareHashAndPasswordWithValidPassword(t *testing.T) {
	password := "some-password"

	hashedPassword, err := BcryptWrapper{}.HashPassword(password)
	require.Nil(t, err)

	fmt.Println(hashedPassword, password)

	err = BcryptWrapper{}.CompareHashAndPassword(hashedPassword, password)
	assert.Nil(t, err)
}