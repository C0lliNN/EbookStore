package usecase

import "github.com/c0llinn/ebook-store/internal/auth/model"

type Repository interface {
	Save(user *model.User) error
	Update(user *model.User) error
	FindByEmail(email string) (model.User, error)
}

type JWTWrapper interface {
	ExtractUserFromToken(tokenString string) (user model.User, err error)
	GenerateTokenForUser(user model.User) (string, error)
}

type EmailClient interface {
	SendPasswordResetEmail(user model.User, newPassword string) error
}

type PasswordGenerator interface {
	NewPassword() string
}