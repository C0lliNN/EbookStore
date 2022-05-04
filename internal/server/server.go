package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Config struct {
	Router                   *gin.Engine
	AuthenticationMiddleware AuthenticationMiddleware
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
	// Errors
	// Swag Handler

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
