package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_IsAdminWhenRoleIsCustomer(t *testing.T) {
	user := User{Role: Customer}
	assert.False(t, user.IsAdmin())
}

func TestUser_IsAdminWhenRoleIsAdmin(t *testing.T) {
	user := User{Role: Admin}

	assert.True(t, user.IsAdmin())
}

func TestUser_FullName(t *testing.T) {
	user := User{FirstName: "Raphael", LastName: "Collin"}

	assert.Equal(t, "Raphael Collin", user.FullName())
}
