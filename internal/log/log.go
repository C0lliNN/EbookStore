package log

// A simple package with initialization

import (
	"context"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
)

var (
	dev  *zap.SugaredLogger
	prod *zap.SugaredLogger
)

func init() {
	d, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	dev = d.Sugar()

	p, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	prod = p.Sugar()
}

// FromContext returns a custom logger based on the context.Context and environment.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	logger := Default()

	logger.With()

	return logger
}

// Default returns a logger based on the environment. It should only be used when context.Context is not available
func Default() *zap.SugaredLogger {
	if strings.EqualFold(viper.GetString("ENV"), "production") {
		return prod
	} else {
		return dev
	}
}
