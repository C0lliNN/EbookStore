//go:build unit
// +build unit

package generator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordGenerator_NewPassword(t *testing.T) {
	password := PasswordGenerator{}.NewPassword()

	assert.Len(t, password, 8)
}
