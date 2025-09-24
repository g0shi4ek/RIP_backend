package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func (h *ChargingHandler) AddTariffToApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	draftApplication, err := h.chargingService.GetDraftChargingApplication(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	chargingOrder, err := h.chargingService.AddChargingOrderToApplication(ctx, tariffId, draftApplication.Id)
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to add tariff to application: %v", err))
		return
	}

	c.JSON(http.StatusCreated, chargingOrder)
}

func (h *ChargingHandler) DeleteChargingOrderFromApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	orderId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	chargingApplication, err := h.chargingService.RemoveChargingOrderFromApplication(ctx, orderId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, chargingApplication)
}

func (h *ChargingHandler) UpdateChargingOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	orderId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var orderRequest struct {
		BatteryCapacity float32 `json:"battery_capacity" binding:"required"`
		CurrentPercent  int     `json:"current_percent" binding:"required"`
		StartTime       string  `json:"start_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		h.ErrorHandler(c, ErrInvalidRequestBody)
		return
	}
	location, _ := time.LoadLocation("Local")

	startTime, err := time.ParseInLocation("15:04", orderRequest.StartTime, location)
	if err != nil {
		h.ErrorHandler(c, ErrInvalidRequestBody)
		return
	}

	order := domain.ChargingOrder{
		Id:              orderId,
		BatteryCapacity: orderRequest.BatteryCapacity,
		CurrentPercent:  orderRequest.CurrentPercent,
		StartTime:       startTime,
	}

	err = h.chargingService.UpdateChargingOrder(ctx, &order)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})
}
