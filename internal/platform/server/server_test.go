package server_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/platform/config"
	"github.com/ebookstore/internal/platform/email"
	"github.com/ebookstore/internal/platform/generator"
	"github.com/ebookstore/internal/platform/hash"
	"github.com/ebookstore/internal/platform/migrator"
	"github.com/ebookstore/internal/platform/payment"
	"github.com/ebookstore/internal/platform/persistence"
	"github.com/ebookstore/internal/platform/server"
	"github.com/ebookstore/internal/platform/storage"
	"github.com/ebookstore/internal/platform/token"
	"github.com/ebookstore/internal/platform/validator"
	"github.com/ebookstore/test"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ServerSuiteTest struct {
	suite.Suite
	baseURL             string
	postgresContainer   *test.PostgresContainer
	localstackContainer *test.LocalstackContainer
	db                  *gorm.DB
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
		_ = s.server.Start()
	}()

	require.Eventually(s.T(), func() bool {
		response, err := http.Get(s.baseURL + "/api/v1/healthcheck")
		if err != nil {
			return false
		}

		return response.StatusCode == http.StatusOK
	}, time.Second*15, time.Millisecond*500)
}

func (s *ServerSuiteTest) TearDownSuite() {
	ctx := context.TODO()

	_ = s.postgresContainer.Terminate(ctx)
	_ = s.localstackContainer.Terminate(ctx)
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
	catalogConfig := catalog.Config{
		Repository:    bookRepository,
		StorageClient: storageStorage,
		IDGenerator:   uuidGenerator,
		Validator:     validatorValidator,
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
		RateLimitMiddleware:      server.NewRateLimitMiddleware(),
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

	s.db = db
	s.server = serverServer
}
