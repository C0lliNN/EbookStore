package container

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebookstore/internal/core/auth"
	"github.com/ebookstore/internal/core/catalog"
	"github.com/ebookstore/internal/core/shop"
	"github.com/ebookstore/internal/log"
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
	"gorm.io/gorm"
)

type Container struct {
	server *server.Server
	dbMigrator *migrator.Migrator
	db    *gorm.DB
}

func New() *Container {
	container := &Container{}

	databaseURI := config.NewMigrationDatabaseURI()
	source := config.NewMigrationSource()
	migratorConfig := migrator.Config{
		DatabaseURI: databaseURI,
		Source:      source,
	}
	dbMigrator := migrator.New(migratorConfig)
	engine := config.NewServerEngine()
	correlationIDMiddleware := server.NewCorrelationIDMiddleware()
	db := config.NewConnection()
	healthcheckHandler := server.NewHeathcheckHandler(db)
	rateLimitMiddleware := server.NewRateLimitMiddleware()
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
	stripePaymentService := payment.NewStripePaymentService()
	shopConfig := shop.Config{
		Repository:     orderRepository,
		PaymentClient:  stripePaymentService,
		CatalogService: catalogCatalog,
		IDGenerator:    uuidGenerator,
		Validator:      validatorValidator,
	}
	shopShop := shop.New(shopConfig)
	shopHandler := server.NewShopHandler(shopShop)
	addr := config.NewServerAddr()
	timeout := config.NewServerTimeout()
	serverConfig := server.Config{
		Router:                   engine,
		CorrelationIDMiddleware:  correlationIDMiddleware,
		HealthcheckHandler:       healthcheckHandler,
		RateLimitMiddleware:      rateLimitMiddleware,
		LoggerMiddleware:         loggerMiddleware,
		AuthenticationMiddleware: authenticationMiddleware,
		ErrorMiddleware:          errorMiddleware,
		AuthenticationHandler:    authenticationHandler,
		CatalogHandler:           catalogHandler,
		ShopHandler:              shopHandler,
		Addr:                     addr,
		Timeout:                  timeout,
	}
	
	container.db = db
	container.dbMigrator = dbMigrator
	container.server = server.New(serverConfig)
	
	return container
}

func (c *Container) Start(ctx context.Context) {
	c.dbMigrator.Sync()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := c.server.Start(); err != nil && err != http.ErrServerClosed {
			log.FromContext(ctx).Fatalf("failed starting server: %v", err)
		}
	}()
	log.FromContext(ctx).Info("starting HTTP server")

	<-done
	log.FromContext(ctx).Info("shutting down HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.server.Shutdown(ctx); err != nil {
		log.FromContext(ctx).Info("failed to shutdown HTTP server")
	}

	log.FromContext(ctx).Info("shutting down database")
	db, err := c.db.DB()
	if err != nil {
		log.FromContext(ctx).Infof("failed to get database instance: %v", err)
	}
    
	if err := db.Close(); err != nil {
		log.FromContext(ctx).Infof("failed to close database: %v", err)
	}

	log.FromContext(ctx).Info("server was shutdown properly")
}

func (c *Container) DB() *gorm.DB {
	return c.db
}