package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type HealthcheckHandler struct{
	db *gorm.DB
}

func NewHeathcheckHandler(db *gorm.DB) *HealthcheckHandler {
	return &HealthcheckHandler{db: db}
}

func (h *HealthcheckHandler) Routes() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/healthcheck", Handler: h.healthcheck, Public: true},
	}
}

// healthcheck godoc
// @Summary REST API Healtcheck
// @Produce  json
// @Success 200 "OK"
// @Failure 500 {object} ErrorResponse
// @Router /healthcheck [get]
func (h *HealthcheckHandler) healthcheck(c *gin.Context) {
	db, err := h.db.DB()
	if err != nil {
		c.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
