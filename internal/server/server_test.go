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
	"gorm.io/gorm"
	"net/http"
	"time"
)

type ServerSuiteTest struct {
	suite.Suite
	baseURL             string
	postgresContainer   *test.PostgresContainer
	localstackContainer *test.LocalstackContainer
	db *gorm.DB
	server              *server.Server
}

func (s *ServerSuiteTest) SetupSuite() {
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
	s.createServer()

	go func() {
		s.server.Start()
	}()

	require.Eventually(s.T(), func() bool {
		response, err := http.Get(s.baseURL + "/healthcheck")
		if err != nil {
			return false
		}

		return response.StatusCode == http.StatusOK
	}, time.Second*15, time.Millisecond*500)
}

func (s *ServerSuiteTest) TearDownSuite() {
	ctx := context.TODO()

	s.postgresContainer.Terminate(ctx)
	s.localstackContainer.Terminate(ctx)
}

func (s *ServerSuiteTest) createServer() {
	databaseURI := config.NewMigrationDatabaseURI()
	source := config.NewMigrationSource()
	migratorConfig := migrator.Config{
		DatabaseURI: databaseURI,
		Source:      source,
	}
	migratorMigrator := migrator.New(migratorConfig)
	engine := config.NewServerEngine()
	hmacSecret := config.NewHMACSecret()
	jwtWrapper := token.NewJWTWrapper(hmacSecret)
	authenticationMiddleware := server.NewAuthenticationMiddleware(jwtWrapper)
	errorMiddleware := server.NewErrorMiddleware()
	db := config.NewConnection()
	healthcheckHandler := server.NewHeathcheckHandler(db)
	userRepository := persistence.NewUserRepository(db)
	bcryptWrapper := hash.NewBcryptWrapper()
	ses := config.NewSESService()
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
		HealthcheckHandler:       healthcheckHandler,
		AuthenticationMiddleware: authenticationMiddleware,
		ErrorMiddleware:          errorMiddleware,
		AuthenticationHandler:    authenticationHandler,
		CatalogHandler:           catalogHandler,
		ShopHandler:              shopHandler,
		Addr:                     addr,
		Timeout:                  timeout,
	}

	s.db = db
	s.server = server.New(serverConfig)
}
