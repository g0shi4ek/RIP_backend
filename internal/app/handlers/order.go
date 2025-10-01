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
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Security BearerAuth
// @Success 200 {object} object "Order deleted successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Order not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingOrders/{id} [delete]
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
        "charging_application": updatedApplication.ToResponse(),
    })
}

// UpdateChargingOrder godoc
// @Summary Update charging order
// @Description Update charging order details in draft application
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body domain.OrderRequest true "Order data"
// @Security BearerAuth
// @Success 200 {object} object "Order updated successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Order not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingOrders/{id} [put]
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
		"charging_order": newChargingOrder.ToResponse(),
	})
}
