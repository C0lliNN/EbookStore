// +build unit

package api

import (
	"fmt"
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromNotValid(t *testing.T) {
	err := &model.ErrNotValid{
		Input: "RegisterRequest",
		Err:   fmt.Errorf("some error"),
	}

	expected := &Error{
		Code:    400,
		Message: "The provided payload is not valid",
		Details: "RegisterRequest not valid: some error",
	}

	actual := fromNotValid(err)

	assert.Equal(t, expected, actual)
}

func TestFromDuplicateKey(t *testing.T) {
	err := &model.ErrDuplicateKey{
		Key: "email",
		Err: fmt.Errorf("some error"),
	}

	expected := &Error{
		Code:    409,
		Message: "this email is already being used",
		Details: "email violation: some error",
	}

	actual := fromDuplicateKey(err)

	assert.Equal(t, expected, actual)
}