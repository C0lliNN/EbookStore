package http

import "github.com/gin-gonic/gin"

func (h AuthHandler) Routes(engine *gin.Engine) {
	engine.POST("/register", h.register)
	engine.POST("/login", h.login)
	engine.POST("/password-reset", h.resetPassword)
}
