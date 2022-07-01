package log

import (
	"context"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	dev  *zap.SugaredLogger
	prod *zap.SugaredLogger
)

func init() {
	devConfig := zap.NewDevelopmentConfig()
	devConfig.DisableCaller = true
	devConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	d, err := devConfig.Build()
	if err != nil {
		panic(err)
	}

	dev = d.Sugar()

	p, err := zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		panic(err)
	}

	prod = p.Sugar()
}

// FromContext returns a custom logger based on the context.Context and environment.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	logger := Default()

	userId := ""
	if id, ok := ctx.Value("userId").(string); ok {
		userId = id
	}

	return logger.With(
		"requestId", ctx.Value("requestId"),
		"userId", userId,
	)
}

// Default returns a logger based on the environment. It should only be used when context.Context is not available
func Default() *zap.SugaredLogger {
	if strings.EqualFold(viper.GetString("ENV"), "production") {
		return prod
	} else {
		return dev
	}
}
