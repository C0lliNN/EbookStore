package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type CorrelationIDMiddleware struct {}

func NewCorrelationIDMiddleware() *CorrelationIDMiddleware {
	return &CorrelationIDMiddleware{}
}

func (*CorrelationIDMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("requestId", fmt.Sprintf("%d", time.Now().UnixNano()))
	}
}
