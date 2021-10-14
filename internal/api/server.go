package api

import (
	_ "github.com/c0llinn/ebook-store/internal/api/docs"
	auth "github.com/c0llinn/ebook-store/internal/auth/delivery/http"
	"github.com/c0llinn/ebook-store/internal/auth/middleware"
	catalog "github.com/c0llinn/ebook-store/internal/catalog/delivery/http"
	shop "github.com/c0llinn/ebook-store/internal/shop/delivery/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

// @host localhost:8081
// @BasePath /
// @query.collection.format multi

func NewHttpServer(router *gin.Engine, authHandler auth.AuthHandler, catalogHandler catalog.CatalogHandler,
	shopHandler shop.ShopHandler, authMiddleware middleware.AuthenticationMiddleware,
	adminMiddleware middleware.AdminMiddleware) *http.Server {

	router.Use(gin.Recovery())
	router.Use(Errors())

	authHandler.Routes(router)

	authorized := router.Group("/")
	authorized.Use(authMiddleware.Handler())

	admin := router.Group("/")
	admin.Use(authMiddleware.Handler(), adminMiddleware.Handler())

	catalogHandler.AuthRoutes(authorized)
	catalogHandler.AdminRoutes(admin)

	unAuthorized := router.Group("/")

	shopHandler.AuthRoutes(authorized)
	shopHandler.UnAuthRoutes(unAuthorized)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func NewRouter() *gin.Engine {
	return gin.New()
}
