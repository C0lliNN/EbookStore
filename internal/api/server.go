package api

import (
	auth "github.com/c0llinn/ebook-store/internal/auth/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func NewRouter() *gin.Engine {
	return gin.New()
}

func NewHttpServer(router *gin.Engine, authHandler auth.AuthHandler) *http.Server {
	router.Use(gin.Recovery())

	authHandler.Routes(router)

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