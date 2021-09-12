package dto

import "github.com/c0llinn/ebook-store/internal/auth/model"

type RegisterRequest struct {
	FirstName            string `json:"firstName" binding:"required,max=150"`
	LastName             string `json:"lastName" binding:"required,max=150"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=6,max=20"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
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

type LoginRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=6,max=20"`
}