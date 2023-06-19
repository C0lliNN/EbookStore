package auth

type RegisterRequest struct {
	FirstName            string `json:"firstName" validate:"required,max=150"`
	LastName             string `json:"lastName" validate:"required,max=150"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=6,max=20"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

func (r RegisterRequest) User(id string) User {
	return User{
		ID:        id,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Role:      Customer,
		Password:  r.Password,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}
