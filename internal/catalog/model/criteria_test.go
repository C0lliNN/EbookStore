// +build unit

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEqualCriteria(t *testing.T) {
	tests := []struct {
		field string
		value interface{}
	} {
		{"title", "Clean Code"},
		{"price", 45.99},
		{"active", true},
	}

	for _, test := range tests {
		expected := Criteria{
			Field:    test.field,
			Operator: "=",
			Value:    test.value,
		}

		actual := NewEqualCriteria(test.field, test.value)

		assert.Equal(t, expected, actual)
	}
}

func TestNewEqualCriteria_WhenValueIsEmpty(t *testing.T) {
	expected := Criteria{}

	actual := NewILikeCriteria("title", "")

	assert.Equal(t, expected, actual)
}

func TestNewILikeCriteria_WhenValueIsNotEmpty(t *testing.T) {
	expected := Criteria{
		Field:    "title",
		Operator: "ILIKE",
		Value:    "%some-value%",
	}

	actual := NewILikeCriteria("title", "some-value")

	assert.Equal(t, expected, actual)
}

func TestCriteria_IsEmpty(t *testing.T) {
	tests := []struct{
		value interface{}
		expected bool
	} {
		{"", true},
		{0, true},
		{"Clean Code", false},
		{45.0, false},
	}

	for _, test := range tests {
		criteria := Criteria{Value: test.value}

		assert.Equal(t, test.expected, criteria.IsEmpty())
	}
}