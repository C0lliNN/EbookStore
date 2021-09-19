package auth

import (
	"github.com/c0llinn/ebook-store/internal/auth/delivery/http"
	"github.com/c0llinn/ebook-store/internal/auth/email"
	"github.com/c0llinn/ebook-store/internal/auth/helper"
	"github.com/c0llinn/ebook-store/internal/auth/repository"
	"github.com/c0llinn/ebook-store/internal/auth/usecase"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	repository.NewUserRepository,
	wire.Bind(new(usecase.Repository), new(repository.UserRepository)),
	helper.NewHMACSecret,
	helper.NewJWTWrapper,
	wire.Bind(new(usecase.JWTWrapper), new(helper.JWTWrapper)),
	usecase.NewAuthUseCase,
	wire.Bind(new(http.UseCase), new(usecase.AuthUseCase)),
	helper.NewUUIDGenerator,
	wire.Bind(new(http.IDGenerator), new(helper.UUIDGenerator)),
	http.NewAuthHandler,
	email.NewEmailClient,
	wire.Bind(new(usecase.EmailClient), new(email.Client)),
	helper.NewPasswordGenerator,
	wire.Bind(new(usecase.PasswordGenerator), new(helper.PasswordGenerator)),
	helper.NewBcryptWrapper,
	wire.Bind(new(usecase.BcryptWrapper), new(helper.BcryptWrapper)),
)
