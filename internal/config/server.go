package config

import (
	"github.com/c0llinn/ebook-store/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

func NewServerEngine() *gin.Engine {
	return gin.New()
}

func NewServerAddr() server.Addr {
	return server.Addr(viper.GetString("SERVER_ADDR"))
}

func NewServerTimeout() server.Timeout {
	return server.Timeout(viper.GetInt64("SERVER_TIMEOUT") * int64(time.Second))
}
