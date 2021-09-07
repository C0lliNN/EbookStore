package usecase

import "github.com/c0llinn/ebook-store/internal/auth/model"

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
