package auth

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

type TokenHandler interface {
	ExtractUserFromToken(tokenString string) (user User, err error)
	GenerateTokenForUser(user User) (string, error)
}

type HashHandler interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type EmailClient interface {
	SendPasswordResetEmail(ctx context.Context, user User, newPassword string) error
}

type PasswordGenerator interface {
	NewPassword() string
}

type IDGenerator interface {
	NewID() string
}

type Config struct {
	Repository        Repository
	Tokener           TokenHandler
	Hasher            HashHandler
	EmailClient       EmailClient
	PasswordGenerator PasswordGenerator
	IDGenerator       IDGenerator
}

type Authenticator struct {
	Config
}

func New(c Config) *Authenticator {
	return &Authenticator{
		Config: c,
	}
}

func (u *Authenticator) Register(ctx context.Context, request RegisterRequest) (CredentialsResponse, error) {
	user := request.User(u.IDGenerator.NewID())

	hashedPassword, err := u.Hasher.HashPassword(user.Password)
	if err != nil {
		return CredentialsResponse{}, err
	}
	user.Password = hashedPassword

	if err = u.Repository.Save(ctx, &user); err != nil {
		return CredentialsResponse{}, err
	}

	return u.generateCredentialsForUser(user)
}

func (u *Authenticator) Login(ctx context.Context, request LoginRequest) (CredentialsResponse, error) {
	user, err := u.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return CredentialsResponse{}, err
	}

	if err = u.Hasher.CompareHashAndPassword(user.Password, request.Password); err != nil {
		return CredentialsResponse{}, ErrWrongPassword
	}

	return u.generateCredentialsForUser(user)
}

func (u *Authenticator) generateCredentialsForUser(user User) (CredentialsResponse, error) {
	token, err := u.Tokener.GenerateTokenForUser(user)
	if err != nil {
		return CredentialsResponse{}, err
	}

	return FromCredentials(Credentials{Token: token}), nil
}

func (u *Authenticator) ResetPassword(ctx context.Context, request PasswordResetRequest) error {
	user, err := u.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	newPassword := u.PasswordGenerator.NewPassword()
	hashedNewPassword, err := u.Hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedNewPassword
	if err = u.Repository.Update(ctx, &user); err != nil {
		return err
	}

	return u.EmailClient.SendPasswordResetEmail(ctx, user, newPassword)
}
