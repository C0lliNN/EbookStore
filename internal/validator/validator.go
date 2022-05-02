package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
