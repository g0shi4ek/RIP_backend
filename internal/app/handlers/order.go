package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// DeleteChargingOrderFromApplication godoc
// @Summary Delete order from application
// @Description Remove charging order from current user's draft application
// @Tags Charging orders
// @Accept json
// @Produce json
// @Param tariffId path int true "Tariff ID"
// @Security BearerAuth
// @Success 200 {object} object "Order deleted successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Order not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/tariffs/{id} [delete]
func (h *ChargingHandler) DeleteChargingOrderFromApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	draftApplication, err := h.chargingService.GetDraftChargingApplicationIfExist(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	tariffId, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	chargingApplication, err := h.chargingService.RemoveChargingOrderFromApplication(ctx, draftApplication.Id, tariffId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "order deleted successfully",
		"charging_application": chargingApplication.ToResponse(),
	})
}

// UpdateChargingOrder godoc
// @Summary Update charging order
// @Description Update charging order details in draft application
// @Tags Charging orders
// @Accept json
// @Produce json
// @Param tariffId path int true "Tariff ID"
// @Param request body domain.ChargingOrderRequest true "Order data"
// @Security BearerAuth
// @Success 200 {object} object "Order updated successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Order not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/tariffs/{id} [put]
func (h *ChargingHandler) UpdateChargingOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	draftApplication, err := h.chargingService.GetDraftChargingApplicationIfExist(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	tariffId, err := helpers.ValidateID(c.Param("id"))
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
		ApplicationId:   draftApplication.Id,
		TariffId:        tariffId,
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
		"message":        "order updated successfully",
		"charging_order": newChargingOrder.ToResponse(),
	})
}

func (h *ChargingHandler) UpdateChargingCalculation(c *gin.Context) {
	// Проверка токена
	authToken := c.GetHeader("X-Auth-Token")
	if authToken != "async123charging" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var request struct {
		ApplicationId   uint    `json:"application_id" binding:"required"`
		TariffId        uint    `json:"tariff_id" binding:"required"`
		EstimatedTime   float32 `json:"estimated_time"`
		CalculatedPrice float32 `json:"calculated_price"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Обновляем расчет заказа
	updates := map[string]interface{}{
		"estimated_time":     request.EstimatedTime,
		"calculated_price":   request.CalculatedPrice,
		"calculation_status": "completed",
	}

	err := h.chargingService.UpdateOrderCalculation(c.Request.Context(), request.ApplicationId, request.TariffId, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	// Обновляем общую цену заявки
	err = h.chargingService.UpdateApplicationTotalPrice(c.Request.Context(), request.ApplicationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update total price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Calculation updated and total price recalculated",
		"application_id": request.ApplicationId,
		"tariff_id":      request.TariffId,
	})
}
