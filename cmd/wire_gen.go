// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from wire.go:

func NewServer() *server.Server {
	databaseURI := config.NewMigrationDatabaseURI()
	source := config.NewMigrationSource()
	migratorConfig := migrator.Config{
		DatabaseURI: databaseURI,
		Source:      source,
	}
	migratorMigrator := migrator.New(migratorConfig)
	engine := config.NewServerEngine()
	correlationIDMiddleware := server.NewCorrelationIDMiddleware()
	db := config.NewConnection()
	healthcheckHandler := server.NewHeathcheckHandler(db)
	loggerMiddleware := server.NewLoggerMiddleware()
	hmacSecret := config.NewHMACSecret()
	jwtWrapper := token.NewJWTWrapper(hmacSecret)
	authenticationMiddleware := server.NewAuthenticationMiddleware(jwtWrapper)
	errorMiddleware := server.NewErrorMiddleware()
	userRepository := persistence.NewUserRepository(db)
	bcryptWrapper := hash.NewBcryptWrapper()
	awsConfig := config.NewAWSConfig()
	client := config.NewSESClient(awsConfig)
	emailEmail := email.NewSESEmailClient(client)
	passwordGenerator := generator.NewPasswordGenerator()
	uuidGenerator := generator.NewUUIDGenerator()
	validatorValidator := validator.New()
	authConfig := auth.Config{
		Repository:        userRepository,
		Tokener:           jwtWrapper,
		Hasher:            bcryptWrapper,
		EmailClient:       emailEmail,
		PasswordGenerator: passwordGenerator,
		IDGenerator:       uuidGenerator,
		Validator:         validatorValidator,
	}
	authenticator := auth.New(authConfig)
	authenticationHandler := server.NewAuthenticatorHandler(authenticator)
	bookRepository := persistence.NewBookRepository(db)
	s3Client := config.NewS3Client(awsConfig)
	presignClient := config.NewPresignClient(s3Client)
	bucket := config.NewBucket()
	storageConfig := storage.Config{
		S3Client:      s3Client,
		PresignClient: presignClient,
		Bucket:        bucket,
	}
	storageStorage := storage.NewStorage(storageConfig)
	filenameGenerator := generator.NewFilenameGenerator()
	catalogConfig := catalog.Config{
		Repository:        bookRepository,
		StorageClient:     storageStorage,
		FilenameGenerator: filenameGenerator,
		IDGenerator:       uuidGenerator,
		Validator:         validatorValidator,
	}
	catalogCatalog := catalog.New(catalogConfig)
	catalogHandler := server.NewCatalogHandler(catalogCatalog)
	orderRepository := persistence.NewOrderRepository(db)
	stripeClient := payment.NewStripePaymentService()
	shopConfig := shop.Config{
		Repository:     orderRepository,
		PaymentClient:  stripeClient,
		CatalogService: catalogCatalog,
		IDGenerator:    uuidGenerator,
		Validator:      validatorValidator,
	}
	shopShop := shop.New(shopConfig)
	shopHandler := server.NewShopHandler(shopShop)
	addr := config.NewServerAddr()
	timeout := config.NewServerTimeout()
	serverConfig := server.Config{
		Migrator:                 migratorMigrator,
		Router:                   engine,
		CorrelationIDMiddleware:  correlationIDMiddleware,
		HealthcheckHandler:       healthcheckHandler,
		LoggerMiddleware:         loggerMiddleware,
		AuthenticationMiddleware: authenticationMiddleware,
		ErrorMiddleware:          errorMiddleware,
		AuthenticationHandler:    authenticationHandler,
		CatalogHandler:           catalogHandler,
		ShopHandler:              shopHandler,
		Addr:                     addr,
		Timeout:                  timeout,
	}
	serverServer := server.New(serverConfig)
	return serverServer
}

// wire.go:

var Set = wire.NewSet(config.NewConnection, generator.NewUUIDGenerator, config.NewSESClient, generator.NewPasswordGenerator, config.NewHMACSecret, persistence.NewUserRepository, email.NewSESEmailClient, token.NewJWTWrapper, hash.NewBcryptWrapper, validator.New, wire.Bind(new(auth.Repository), new(*persistence.UserRepository)), wire.Bind(new(auth.Validator), new(*validator.Validator)), wire.Bind(new(auth.TokenHandler), new(*token.JWTWrapper)), wire.Bind(new(auth.HashHandler), new(*hash.BcryptWrapper)), wire.Bind(new(auth.EmailClient), new(*email.Email)), wire.Bind(new(auth.IDGenerator), new(*generator.UUIDGenerator)), wire.Bind(new(auth.PasswordGenerator), new(*generator.PasswordGenerator)), wire.NewSet(wire.Struct(new(auth.Config), "*")), auth.New, config.NewAWSConfig, config.NewBucket, config.NewS3Client, config.NewPresignClient, wire.NewSet(wire.Struct(new(storage.Config), "*")), storage.NewStorage, generator.NewFilenameGenerator, persistence.NewBookRepository, wire.Bind(new(catalog.Repository), new(*persistence.BookRepository)), wire.Bind(new(catalog.Validator), new(*validator.Validator)), wire.Bind(new(catalog.FilenameGenerator), new(*generator.FilenameGenerator)), wire.Bind(new(catalog.IDGenerator), new(*generator.UUIDGenerator)), wire.Bind(new(catalog.StorageClient), new(*storage.Storage)), wire.NewSet(wire.Struct(new(catalog.Config), "*")), catalog.New, persistence.NewOrderRepository, payment.NewStripePaymentService, wire.Bind(new(shop.Repository), new(*persistence.OrderRepository)), wire.Bind(new(shop.Validator), new(*validator.Validator)), wire.Bind(new(shop.PaymentClient), new(*payment.StripePaymentService)), wire.Bind(new(shop.CatalogService), new(*catalog.Catalog)), wire.Bind(new(shop.IDGenerator), new(*generator.UUIDGenerator)), wire.NewSet(wire.Struct(new(shop.Config), "*")), shop.New, config.NewMigrationDatabaseURI, config.NewMigrationSource, wire.NewSet(wire.Struct(new(migrator.Config), "*")), migrator.New, config.NewServerEngine, config.NewServerAddr, config.NewServerTimeout, server.NewCorrelationIDMiddleware, server.NewErrorMiddleware, server.NewHeathcheckHandler, server.NewLoggerMiddleware, server.NewAuthenticationMiddleware, server.NewAuthenticatorHandler, server.NewCatalogHandler, server.NewShopHandler, wire.Bind(new(server.Authenticator), new(*auth.Authenticator)), wire.Bind(new(server.Catalog), new(*catalog.Catalog)), wire.Bind(new(server.Shop), new(*shop.Shop)), wire.Bind(new(server.TokenHandler), new(*token.JWTWrapper)), wire.NewSet(wire.Struct(new(server.Config), "*")), server.New)
