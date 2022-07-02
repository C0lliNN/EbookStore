package config

import (
	"github.com/ebookstore/internal/token"
	"github.com/spf13/viper"
)

func NewHMACSecret() token.HMACSecret {
	return []byte(viper.GetString("JWT_SECRET"))
}
