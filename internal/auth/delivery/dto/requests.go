package dto

import "github.com/c0llinn/ebook-store/internal/auth/model"

type RegisterRequest struct {
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (r RegisterRequest) ToDomain(id string) model.User {
	return model.User{
		ID:        id,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Role:      model.Customer,
		Password:  r.Password,
	}
}
