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

func (a *Authenticator) Register(ctx context.Context, request RegisterRequest) (CredentialsResponse, error) {
	user := request.User(a.IDGenerator.NewID())

	hashedPassword, err := a.Hasher.HashPassword(user.Password)
	if err != nil {
		return CredentialsResponse{}, err
	}
	user.Password = hashedPassword

	if err = a.Repository.Save(ctx, &user); err != nil {
		return CredentialsResponse{}, err
	}

	return a.generateCredentialsForUser(user)
}

func (a *Authenticator) Login(ctx context.Context, request LoginRequest) (CredentialsResponse, error) {
	user, err := a.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return CredentialsResponse{}, err
	}

	if err = a.Hasher.CompareHashAndPassword(user.Password, request.Password); err != nil {
		return CredentialsResponse{}, ErrWrongPassword
	}

	return a.generateCredentialsForUser(user)
}

func (a *Authenticator) generateCredentialsForUser(user User) (CredentialsResponse, error) {
	token, err := a.Tokener.GenerateTokenForUser(user)
	if err != nil {
		return CredentialsResponse{}, err
	}

	return FromCredentials(Credentials{Token: token}), nil
}

func (a *Authenticator) ResetPassword(ctx context.Context, request PasswordResetRequest) error {
	user, err := a.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	newPassword := a.PasswordGenerator.NewPassword()
	hashedNewPassword, err := a.Hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedNewPassword
	if err = a.Repository.Update(ctx, &user); err != nil {
		return err
	}

	return a.EmailClient.SendPasswordResetEmail(ctx, user, newPassword)
}
