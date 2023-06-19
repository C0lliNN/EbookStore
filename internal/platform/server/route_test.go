package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_IsPublic(t *testing.T) {
	tests := []struct {
		Public   bool
		Expected bool
	}{
		{false, false},
		{true, true},
	}

	for _, tc := range tests {
		route := Route{Public: tc.Public}
		assert.Equal(t, tc.Expected, route.IsPublic())
	}
}
