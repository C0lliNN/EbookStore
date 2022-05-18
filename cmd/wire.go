//go:build wireinject
// +build wireinject

package main

import (
	"github.com/c0llinn/ebook-store/internal/auth"
	"github.com/c0llinn/ebook-store/internal/catalog"
	"github.com/c0llinn/ebook-store/internal/config"
	"github.com/c0llinn/ebook-store/internal/email"
	"github.com/c0llinn/ebook-store/internal/generator"
	"github.com/c0llinn/ebook-store/internal/hash"
	"github.com/c0llinn/ebook-store/internal/migrator"
	"github.com/c0llinn/ebook-store/internal/payment"
	"github.com/c0llinn/ebook-store/internal/persistence"
	"github.com/c0llinn/ebook-store/internal/server"
	"github.com/c0llinn/ebook-store/internal/shop"
	"github.com/c0llinn/ebook-store/internal/storage"
	"github.com/c0llinn/ebook-store/internal/token"
	"github.com/c0llinn/ebook-store/internal/validator"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	config.NewConnection,
	generator.NewUUIDGenerator,
	config.NewSNSService,
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
	wire.Bind(new(auth.EmailClient), new(*email.SESEmailClient)),
	wire.Bind(new(auth.IDGenerator), new(*generator.UUIDGenerator)),
	wire.Bind(new(auth.PasswordGenerator), new(*generator.PasswordGenerator)),
	wire.NewSet(wire.Struct(new(auth.Config), "*")),
	auth.New,

	config.NewBucket,
	config.NewS3Service,
	storage.NewS3Client,
	generator.NewFilenameGenerator,
	persistence.NewBookRepository,
	wire.Bind(new(catalog.Repository), new(*persistence.BookRepository)),
	wire.Bind(new(catalog.Validator), new(*validator.Validator)),
	wire.Bind(new(catalog.FilenameGenerator), new(*generator.FilenameGenerator)),
	wire.Bind(new(catalog.IDGenerator), new(*generator.UUIDGenerator)),
	wire.Bind(new(catalog.StorageClient), new(*storage.S3Client)),
	wire.NewSet(wire.Struct(new(catalog.Config), "*")),
	catalog.New,

	persistence.NewOrderRepository,
	payment.NewStripeClient,
	wire.Bind(new(shop.Repository), new(*persistence.OrderRepository)),
	wire.Bind(new(shop.Validator), new(*validator.Validator)),
	wire.Bind(new(shop.PaymentClient), new(*payment.StripeClient)),
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
	server.NewErrorMiddleware,
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
