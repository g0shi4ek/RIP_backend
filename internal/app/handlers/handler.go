package handlers

import (
	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChargingHandler struct {
	chargingRepository domain.IChargingRepository
}

func NewChargingHandler(repo domain.IChargingRepository) (*ChargingHandler, error) {
	return &ChargingHandler{
		chargingRepository: repo,
	}, nil
}

func (h *ChargingHandler) RegisterChargingHandler(r *gin.Engine) {
	r.GET("/tariffs", h.GetTarrifs)
	r.GET("/tariff/:id", h.GetChargingTarrifById)
	r.GET("/application", h.GetChargingApplicationByStatus)
	r.POST("/application/:id", h.DeleteChargingApplicationById)
	r.POST("/tariff/:tariffId/application", h.AddTariffToApplication)
}

func (h *ChargingHandler) RegisterChargingStatic(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/resources", "./resources")
}

func (h *ChargingHandler) ErrorChargingHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
