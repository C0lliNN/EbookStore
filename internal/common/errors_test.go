// +build unit

package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrNotValid_Error(t *testing.T) {
	err := ErrNotValid{
		Input: "RegisterRequest",
		Err:   fmt.Errorf("some error"),
	}

	assert.Equal(t, "RegisterRequest not valid: some error", err.Error())
}

func TestErrDuplicateKey_Error(t *testing.T) {
	err := ErrDuplicateKey{
		Key: "email",
		Err: fmt.Errorf("some error"),
	}

	assert.Equal(t, "email violation: some error", err.Error())
}

func TestErrEntityNotFound_Error(t *testing.T) {
	err := ErrEntityNotFound{
		Entity: "User",
		Err:    fmt.Errorf("some error"),
	}

	assert.Equal(t, "User could not be found: some error", err.Error())
}

func TestErrWrongPassword_Error(t *testing.T) {
	err := ErrWrongPassword{
		Err:    fmt.Errorf("some error"),
	}

	assert.Equal(t, "some error", err.Error())
}