package usecase

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Save(user *model.User) error
	FindByEmail(email string) (model.User, error)
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

func (u AuthUseCase) Register(user model.User) (model.Credentials, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return model.Credentials{}, err
	}

	user.Password = string(hashedPassword)
	if err = u.repo.Save(&user); err != nil {
		return model.Credentials{}, err
	}

	return u.generateCredentialsForUser(user)
}

func (u AuthUseCase) Login(email, password string) (model.Credentials, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return model.Credentials{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.Credentials{}, &model.ErrWrongPassword{Err: err}
	}

	return u.generateCredentialsForUser(user)
}

func (u AuthUseCase) generateCredentialsForUser(user model.User) (model.Credentials, error) {
	token, err := u.jwt.GenerateTokenForUser(user)
	if err != nil {
		return model.Credentials{}, err
	}

	return model.Credentials{Token: token}, nil
}