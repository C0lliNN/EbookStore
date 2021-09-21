package usecase

import (
	"github.com/c0llinn/ebook-store/internal/auth/model"
	"github.com/c0llinn/ebook-store/internal/common"
)

type AuthUseCase struct {
	repo              Repository
	jwt               JWTWrapper
	bcrypt            BcryptWrapper
	emailClient       EmailClient
	passwordGenerator PasswordGenerator
}

func NewAuthUseCase(repo Repository, jwt JWTWrapper, emailClient EmailClient, passwordGenerator PasswordGenerator, bcryptWrapper BcryptWrapper) AuthUseCase {
	return AuthUseCase{
		repo:              repo,
		jwt:               jwt,
		bcrypt:            bcryptWrapper,
		emailClient:       emailClient,
		passwordGenerator: passwordGenerator,
	}
}

func (u AuthUseCase) Register(user model.User) (model.Credentials, error) {
	hashedPassword, err := u.bcrypt.HashPassword(user.Password)
	if err != nil {
		return model.Credentials{}, err
	}

	user.Password = hashedPassword
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

	if err = u.bcrypt.CompareHashAndPassword(user.Password, password); err != nil {
		return model.Credentials{}, &common.ErrWrongPassword{Err: err}
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

func (u AuthUseCase) ResetPassword(email string) error {
	user, err := u.repo.FindByEmail(email)
	if err != nil {
		return err
	}

	newPassword := u.passwordGenerator.NewPassword()
	hashedNewPassword, err := u.bcrypt.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedNewPassword
	if err = u.repo.Update(&user); err != nil {
		return err
	}

	return u.emailClient.SendPasswordResetEmail(user, newPassword)
}
