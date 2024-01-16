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

// fromContext returns a custom logger based on the context.Context and environment.
func fromContext(ctx context.Context) *zap.SugaredLogger {
	logger := dev
	if strings.EqualFold(viper.GetString("ENV"), "production") {
		logger = prod
	}

	userId := ""
	if id, ok := ctx.Value("userId").(string); ok {
		userId = id
	}

	return logger.With(
		"requestId", ctx.Value("requestId"),
		"userId", userId,
	)
}

func With(ctx context.Context, args ...interface{}) *zap.SugaredLogger {
	return fromContext(ctx).With(args...)
}

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	fromContext(ctx).Debugf(msg, args...)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	fromContext(ctx).Infof(msg, args...)
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	fromContext(ctx).Warnf(msg, args...)
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	fromContext(ctx).Errorf(msg, args...)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	fromContext(ctx).Fatalf(msg, args...)
}