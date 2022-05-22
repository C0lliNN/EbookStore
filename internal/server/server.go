package server

import (
	_ "github.com/c0llinn/ebook-store/docs"
	"github.com/c0llinn/ebook-store/internal/migrator"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
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
	HealthcheckHandler *HealthcheckHandler
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
}

func New(c Config) *Server {
	return &Server{Config: c}
}

func (s *Server) Start() error {
	s.Migrator.Sync()

	router := s.Router

	router.Use(gin.Recovery())
	router.Use(s.ErrorMiddleware.Handler())

	// This is redirect is for convenience purposes. It's easy to remember /docs
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := s.AuthenticationHandler.Routes()
	routes = append(routes, s.HealthcheckHandler.Routes()...)
	routes = append(routes, s.CatalogHandler.Routes()...)
	routes = append(routes, s.ShopHandler.Routes()...)

	authorizedRouter := router.Group("/", s.AuthenticationMiddleware.Handler())

	for _, r := range routes {
		if r.IsPublic() {
			router.Handle(r.Method, r.Path, r.Handler)
		} else {
			authorizedRouter.Handle(r.Method, r.Path, r.Handler)
		}
	}

	httpServer := &http.Server{
		Handler:      router,
		Addr:         string(s.Addr),
		WriteTimeout: time.Duration(s.Timeout),
		ReadTimeout:  time.Duration(s.Timeout),
	}

	return httpServer.ListenAndServe()
}
