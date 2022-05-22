package server_test

import (
	"context"
	"fmt"
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
	"github.com/c0llinn/ebook-store/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"time"
)

type ServerTest struct {
	suite.Suite
	baseURL             string
	postgresContainer   *test.PostgresContainer
	localstackContainer *test.LocalstackContainer
	server              *server.Server
}

func (s *ServerTest) SetupSuite() {
	config.LoadConfiguration()

	var err error
	ctx := context.TODO()

	s.postgresContainer, err = test.NewPostgresContainer(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.localstackContainer, err = test.NewLocalstackContainer(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	viper.Set("DATABASE_URI", s.postgresContainer.URI)
	viper.Set("AWS_SES_ENDPOINT", fmt.Sprintf("http://localhost:%v", s.localstackContainer.Port))
	viper.Set("AWS_S3_ENDPOINT", fmt.Sprintf("http://s3.localhost.localstack.cloud:%v", s.localstackContainer.Port))

	s.baseURL = fmt.Sprintf("http://%v", viper.GetString("SERVER_ADDR"))
	s.server = newServer()

	go func() {
		s.server.Start()
	}()

	require.Eventually(s.T(), func() bool {
		s.server.
	}, time.Second*2, time.Millisecond*100)
}

func (s *ServerTest) TearDownSuite() {
	ctx := context.TODO()

	s.postgresContainer.Terminate(ctx)
	s.localstackContainer.Terminate(ctx)
}

func newServer() *server.Server {
	databaseURL := config.NewMigrationDatabaseURI()
	migrationSource := config.NewMigrationSource()
	migratorConfig := migrator.Config{
		DatabaseURI: databaseURL,
		Source:      migrationSource,
	}
	migratorMigrator := migrator.New(migratorConfig)
	engine := config.NewServerEngine()
	hmacSecret := config.NewHMACSecret()
	jwtWrapper := token.NewJWTWrapper(hmacSecret)
	authenticationMiddleware := server.NewAuthenticationMiddleware(jwtWrapper)
	errorMiddleware := server.NewErrorMiddleware()
	db := config.NewConnection()
	userRepository := persistence.NewUserRepository(db)
	bcryptWrapper := hash.NewBcryptWrapper()
	ses := config.NewSNSService()
	sesEmailClient := email.NewSESEmailClient(ses)
	passwordGenerator := generator.NewPasswordGenerator()
	uuidGenerator := generator.NewUUIDGenerator()
	validatorValidator := validator.New()
	authConfig := auth.Config{
		Repository:        userRepository,
		Tokener:           jwtWrapper,
		Hasher:            bcryptWrapper,
		EmailClient:       sesEmailClient,
		PasswordGenerator: passwordGenerator,
		IDGenerator:       uuidGenerator,
		Validator:         validatorValidator,
	}
	authenticator := auth.New(authConfig)
	authenticationHandler := server.NewAuthenticatorHandler(authenticator)
	bookRepository := persistence.NewBookRepository(db)
	s3 := config.NewS3Service()
	bucket := config.NewBucket()
	s3Client := storage.NewS3Client(s3, bucket)
	filenameGenerator := generator.NewFilenameGenerator()
	catalogConfig := catalog.Config{
		Repository:        bookRepository,
		StorageClient:     s3Client,
		FilenameGenerator: filenameGenerator,
		IDGenerator:       uuidGenerator,
		Validator:         validatorValidator,
	}
	catalogCatalog := catalog.New(catalogConfig)
	catalogHandler := server.NewCatalogHandler(catalogCatalog)
	orderRepository := persistence.NewOrderRepository(db)
	stripeClient := payment.NewStripeClient()
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
		AuthenticationMiddleware: authenticationMiddleware,
		ErrorMiddleware:          errorMiddleware,
		AuthenticationHandler:    authenticationHandler,
		CatalogHandler:           catalogHandler,
		ShopHandler:              shopHandler,
		Addr:                     addr,
		Timeout:                  timeout,
	}

	return server.New(serverConfig)
}
