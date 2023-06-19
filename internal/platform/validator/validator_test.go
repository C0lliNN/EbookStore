package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_Validate(t *testing.T) {
	type person struct {
		Name string `validate:"required"`
	}

	tests := []struct {
		Name          string
		Person        person
		ExpectedError bool
	}{
		{
			Name:          "Invalid input",
			Person:        person{Name: ""},
			ExpectedError: true,
		},
		{
			Name:          "Valid input",
			Person:        person{Name: "Raphael"},
			ExpectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			err := New().Validate(tc.Person)
			assert.Equal(t, tc.ExpectedError, err != nil)
		})
	}

}
