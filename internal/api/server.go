package api

import (
	auth "github.com/c0llinn/ebook-store/internal/auth/delivery/http"
	"github.com/c0llinn/ebook-store/internal/auth/middleware"
	catalog "github.com/c0llinn/ebook-store/internal/catalog/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func NewRouter() *gin.Engine {
	return gin.New()
}

func NewHttpServer(router *gin.Engine, authHandler auth.AuthHandler, catalogHandler catalog.CatalogHandler,
	authMiddleware middleware.AuthenticationMiddleware) *http.Server {
	router.Use(gin.Recovery())
	router.Use(Errors())

	authHandler.Routes(router)

	authorized := router.Group("/")
	authorized.Use(authMiddleware.Handler())

	catalogHandler.Routes(authorized)

	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	return &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}