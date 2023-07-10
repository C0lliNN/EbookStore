package server

import (
	"context"
	"net/http"
	"time"

	_ "github.com/ebookstore/docs"
	"github.com/ebookstore/internal/platform/migrator"
	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title E-book Store
// @version 1.0
// @description Endpoints available in the E-book store REST API.
// @termsOfService https://github.com/C0lliNN

// @contact.name Raphael Collin
// @contact.email raphael_professional@yahoo.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host http://localhost:8080
// @BasePath /
// @query.collection.format multi

type Addr string
type Timeout time.Duration

type Config struct {
	Migrator                 *migrator.Migrator
	Router                   *gin.Engine
	CorrelationIDMiddleware  *CorrelationIDMiddleware
	HealthcheckHandler       *HealthcheckHandler
	RateLimitMiddleware      *RateLimitMiddleware
	LoggerMiddleware         *LoggerMiddleware
	AuthenticationMiddleware *AuthenticationMiddleware
	ErrorMiddleware          *ErrorMiddleware
	AuthenticationHandler    *AuthenticationHandler
	CatalogHandler           *CatalogHandler
	ShopHandler              *ShopHandler
	Addr                     Addr
	Timeout                  Timeout
}

type Server struct {
	Config
	httpServer *http.Server
}

func New(c Config) *Server {
	return &Server{Config: c}
}

func (s *Server) Start() error {
	s.Migrator.Sync()

	router := s.Router

	router.Use(s.CorrelationIDMiddleware.Handler())
	router.Use(gin.Recovery())
	router.Use(limits.RequestSizeLimiter(2 << 20)) // 10MB

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	router.Use(cors.New(corsConfig))
	router.Use(s.RateLimitMiddleware.Handler())
	router.Use(s.LoggerMiddleware.Handler())
	router.Use(s.ErrorMiddleware.Handler())

	// This redirect is for convenience purposes. It's easy to remember /docs
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := s.AuthenticationHandler.Routes()
	routes = append(routes, s.HealthcheckHandler.Routes()...)
	routes = append(routes, s.CatalogHandler.Routes()...)
	routes = append(routes, s.ShopHandler.Routes()...)

	versionedRouter := router.Group("/api/v1")
	authorizedRouter := versionedRouter.Group("/", s.AuthenticationMiddleware.Handler())

	for _, r := range routes {
		if r.IsPublic() {
			versionedRouter.Handle(r.Method, r.Path, r.Handler)
		} else {
			authorizedRouter.Handle(r.Method, r.Path, r.Handler)
		}
	}

	s.httpServer = &http.Server{
		Handler:      router,
		Addr:         string(s.Addr),
		WriteTimeout: time.Duration(s.Timeout),
		ReadTimeout:  time.Duration(s.Timeout),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
