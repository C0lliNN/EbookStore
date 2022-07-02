//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ebookstore/internal/auth"
	"github.com/ebookstore/internal/catalog"
	"github.com/ebookstore/internal/config"
	"github.com/ebookstore/internal/email"
	"github.com/ebookstore/internal/generator"
	"github.com/ebookstore/internal/hash"
	"github.com/ebookstore/internal/migrator"
	"github.com/ebookstore/internal/payment"
	"github.com/ebookstore/internal/persistence"
	"github.com/ebookstore/internal/server"
	"github.com/ebookstore/internal/shop"
	"github.com/ebookstore/internal/storage"
	"github.com/ebookstore/internal/token"
	"github.com/ebookstore/internal/validator"
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

	generator.NewFilenameGenerator,
	persistence.NewBookRepository,
	wire.Bind(new(catalog.Repository), new(*persistence.BookRepository)),
	wire.Bind(new(catalog.Validator), new(*validator.Validator)),
	wire.Bind(new(catalog.FilenameGenerator), new(*generator.FilenameGenerator)),
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
