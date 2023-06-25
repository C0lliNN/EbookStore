//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/migrator"
	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/internal/platform/email"
	"github.com/ebookstore/internal/platform/generator"
	"github.com/ebookstore/internal/platform/hash"
	"github.com/ebookstore/internal/platform/payment"
	"github.com/ebookstore/internal/platform/persistence"
	"github.com/ebookstore/internal/platform/server"
	"github.com/ebookstore/internal/platform/storage"
	"github.com/ebookstore/internal/platform/token"
	"github.com/ebookstore/internal/platform/validator"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	config.NewConnection,
	generator.NewUUIDGenerator,
	config.NewSESClient,
	generator.NewPasswordGenerator,
	config.NewHMACSecret,
	persistence.NewUserRepository,
	email.NewSESEmailClient,
	token.NewJWTWrapper,
	hash.NewBcryptWrapper,
	validator.New,
	wire.Bind(new(auth.Repository), new(*persistence.UserRepository)),
	wire.Bind(new(auth.Validator), new(*validator.Validator)),
	wire.Bind(new(auth.TokenHandler), new(*token.JWTWrapper)),
	wire.Bind(new(auth.HashHandler), new(*hash.BcryptWrapper)),
	wire.Bind(new(auth.EmailClient), new(*email.Email)),
	wire.Bind(new(auth.IDGenerator), new(*generator.UUIDGenerator)),
	wire.Bind(new(auth.PasswordGenerator), new(*generator.PasswordGenerator)),
	wire.NewSet(wire.Struct(new(auth.Config), "*")),
	auth.New,

	config.NewAWSConfig,
	config.NewBucket,
	config.NewS3Client,
	config.NewPresignClient,
	wire.NewSet(wire.Struct(new(storage.Config), "*")),
	storage.NewStorage,

	persistence.NewBookRepository,
	wire.Bind(new(catalog.Repository), new(*persistence.BookRepository)),
	wire.Bind(new(catalog.Validator), new(*validator.Validator)),
	wire.Bind(new(catalog.IDGenerator), new(*generator.UUIDGenerator)),
	wire.Bind(new(catalog.StorageClient), new(*storage.Storage)),
	wire.NewSet(wire.Struct(new(catalog.Config), "*")),
	catalog.New,

	persistence.NewOrderRepository,
	payment.NewStripePaymentService,
	wire.Bind(new(shop.Repository), new(*persistence.OrderRepository)),
	wire.Bind(new(shop.Validator), new(*validator.Validator)),
	wire.Bind(new(shop.PaymentClient), new(*payment.StripePaymentService)),
	wire.Bind(new(shop.CatalogService), new(*catalog.Catalog)),
	wire.Bind(new(shop.IDGenerator), new(*generator.UUIDGenerator)),
	wire.NewSet(wire.Struct(new(shop.Config), "*")),
	shop.New,

	config.NewMigrationDatabaseURI,
	config.NewMigrationSource,
	wire.NewSet(wire.Struct(new(migrator.Config), "*")),
	migrator.New,

	config.NewServerEngine,
	config.NewServerAddr,
	config.NewServerTimeout,
	server.NewCorrelationIDMiddleware,
	server.NewErrorMiddleware,
	server.NewHeathcheckHandler,
	server.NewRateLimitMiddleware,
	server.NewLoggerMiddleware,
	server.NewAuthenticationMiddleware,
	server.NewAuthenticatorHandler,
	server.NewCatalogHandler,
	server.NewShopHandler,
	wire.Bind(new(server.Authenticator), new(*auth.Authenticator)),
	wire.Bind(new(server.Catalog), new(*catalog.Catalog)),
	wire.Bind(new(server.Shop), new(*shop.Shop)),
	wire.Bind(new(server.TokenHandler), new(*token.JWTWrapper)),
	wire.NewSet(wire.Struct(new(server.Config), "*")),
	server.New,
)

func NewServer() *server.Server {
	wire.Build(Set)

	return &server.Server{}
}
