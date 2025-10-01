package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func (h *ChargingHandler) DeleteChargingOrderFromApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	chargingOrderId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	updatedApplication, err := h.chargingService.RemoveChargingOrderFromApplication(ctx, chargingOrderId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
        "message":     "order deleted successfully",
        "application": updatedApplication.ToResponse(),
    })
}

func (h *ChargingHandler) UpdateChargingOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	chargingOrderId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var chargingOrderRequest domain.ChargingOrderRequest
	if err := c.ShouldBindJSON(&chargingOrderRequest); err != nil {
		h.ErrorHandler(c, ErrInvalidRequestBody)
		return
	}
	location, _ := time.LoadLocation("Local")

	startTime, err := time.ParseInLocation("15:04", chargingOrderRequest.StartTime, location)
	if err != nil {
		h.ErrorHandler(c, ErrInvalidRequestBody)
		return
	}

	chargingOrder := domain.ChargingOrder{
		Id:              chargingOrderId,
		BatteryCapacity: chargingOrderRequest.BatteryCapacity,
		CurrentPercent:  chargingOrderRequest.CurrentPercent,
		StartTime:       startTime,
	}

	newChargingOrder, err := h.chargingService.UpdateChargingOrder(ctx, &chargingOrder)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order updated successfully",
		"order": newChargingOrder.ToResponse(),
	})
}
