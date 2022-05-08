package server

import (
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

// @host https://ebook-store2.herokuapp.com
// @BasePath /
// @query.collection.format multi

type Config struct {
	Router                   *gin.Engine
	AuthenticationMiddleware AuthenticationMiddleware
	ErrorMiddleware          ErrorMiddleware
	AuthenticationHandler    AuthenticationHandler
	CatalogHandler           CatalogHandler
	ShopHandler              ShopHandler
	Addr                     string
	Timeout                  time.Duration
}

type Server struct {
	Config
}

func New(c Config) *Server {
	return &Server{Config: c}
}

func (s *Server) Start() error {
	router := s.Router

	router.Use(gin.Recovery())
	router.Use(s.ErrorMiddleware.Handler())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := s.AuthenticationHandler.Routes()
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
		Addr:         s.Addr,
		WriteTimeout: s.Timeout,
		ReadTimeout:  s.Timeout,
	}

	return httpServer.ListenAndServe()
}
