package auth

import (
	"context"
	"fmt"

	"github.com/ebookstore/internal/log"
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

type Validator interface {
	Validate(i interface{}) error
}

type Config struct {
	Repository        Repository
	Tokener           TokenHandler
	Hasher            HashHandler
	EmailClient       EmailClient
	PasswordGenerator PasswordGenerator
	IDGenerator       IDGenerator
	Validator         Validator
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
	if err := a.Validator.Validate(request); err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Register) failed validating request: %w", err)
	}

	user := request.User(a.IDGenerator.NewID())

	log.FromContext(ctx).Infof("creating new user with id %s", user.ID)

	hashedPassword, err := a.Hasher.HashPassword(user.Password)
	if err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Register) failed hashing password: %w", err)
	}
	user.Password = hashedPassword

	if err = a.Repository.Save(ctx, &user); err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Register) failed saving user: %w", err)
	}

	credentials, err := a.generateCredentialsForUser(user)
	if err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Register) failed generating credentials: %w", err)
	}

	return credentials, err
}

func (a *Authenticator) Login(ctx context.Context, request LoginRequest) (CredentialsResponse, error) {
	if err := a.Validator.Validate(request); err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Login) failed validating request: %w", err)
	}

	user, err := a.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Login) failed finding user: %w", err)
	}

	log.FromContext(ctx).Info("new login attempt for user with id %s", user.ID)

	if err = a.Hasher.CompareHashAndPassword(user.Password, request.Password); err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Login) failed comparing hash and password: %w", ErrWrongPassword)
	}

	credentials, err := a.generateCredentialsForUser(user)
	if err != nil {
		return CredentialsResponse{}, fmt.Errorf("(Login) failed generating credentials: %w", err)
	}

	return credentials, err
}

func (a *Authenticator) generateCredentialsForUser(user User) (CredentialsResponse, error) {
	token, err := a.Tokener.GenerateTokenForUser(user)
	if err != nil {
		return CredentialsResponse{}, fmt.Errorf("(generateCredentialsForUser) failed generating token: %w", err)
	}

	return NewCredentialsResponse(Credentials{Token: token}), nil
}

func (a *Authenticator) ResetPassword(ctx context.Context, request PasswordResetRequest) error {
	if err := a.Validator.Validate(request); err != nil {
		return fmt.Errorf("(ResetPassword) failed validating request: %w", err)
	}

	user, err := a.Repository.FindByEmail(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("(ResetPassword) failed finding user: %w", err)
	}

	log.FromContext(ctx).Infof("resetting password for user with id %s", user.ID)

	newPassword := a.PasswordGenerator.NewPassword()
	hashedNewPassword, err := a.Hasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("(ResetPassword) failed hashing password: %w", err)
	}

	user.Password = hashedNewPassword
	if err = a.Repository.Update(ctx, &user); err != nil {
		return fmt.Errorf("(ResetPassword) failed updating user: %w", err)
	}

	if err = a.EmailClient.SendPasswordResetEmail(ctx, user, newPassword); err != nil {
		return fmt.Errorf("(ResetPassword) failed sending email: %w", err)
	}

	return nil
}
