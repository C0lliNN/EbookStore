package usecase

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Save(user *model.User) error
}

type JWTWrapper interface {
	ExtractUserFromToken(tokenString string) (user model.User, err error)
	GenerateTokenForUser(user model.User) (string, error)
}

type AuthUseCase struct {
	repo Repository
	jwt JWTWrapper
}

func NewAuthUseCase(repo Repository, jwt JWTWrapper) AuthUseCase {
	return AuthUseCase{repo: repo, jwt: jwt}
}

func (u AuthUseCase) Register(user model.User) (credentials model.Credentials, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return
	}

	user.Password = string(hashedPassword)
	if err = u.repo.Save(&user); err != nil {
		return
	}

	token, err := u.jwt.GenerateTokenForUser(user)
	if err != nil {
		return
	}

	credentials.Token = token
	return
}
