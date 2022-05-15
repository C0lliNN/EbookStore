package config

import (
	"github.com/c0llinn/ebook-store/internal/token"
	"github.com/spf13/viper"
)

func NewHMACSecret() token.HMACSecret {
	return []byte(viper.GetString("JWT_SECRET"))
}
