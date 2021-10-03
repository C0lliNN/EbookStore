// +build unit

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCriteria_IsEmpty(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{value: nil, expected: true},
		{value: "", expected: true},
		{value: " ", expected: false},
	}

	for _, test := range tests {
		criteria := Criteria{Field: "test", Operator: "=", Value: test.value}
		assert.Equal(t, test.expected, criteria.IsEmpty())
	}
}

func TestNewEqualCriteria(t *testing.T) {
	expected := Criteria{
		Field:    "user_id",
		Operator: "=",
		Value:    "some-id",
	}

	actual := NewEqualCriteria("user_id", "some-id")

	assert.Equal(t, expected, actual)
}
