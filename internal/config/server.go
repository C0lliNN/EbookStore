package config

import (
	"time"

	"github.com/ebookstore/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
