package usecase

import (
	"context"
	"github.com/c0llinn/ebook-store/internal/auth/model"
)

type Repository interface {
	Save(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

type JWTWrapper interface {
	ExtractUserFromToken(tokenString string) (user model.User, err error)
	GenerateTokenForUser(user model.User) (string, error)
}

type BcryptWrapper interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type EmailClient interface {
	SendPasswordResetEmail(ctx context.Context, user model.User, newPassword string) error
}

type PasswordGenerator interface {
	NewPassword() string
}
