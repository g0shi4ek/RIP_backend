package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *ChargingHandler) AddTariffToApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffId, err := h.validateID(c.Param("id"))
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

	orderId, err := h.validateID(c.Param("id"))
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

	id, err := h.validateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var order domain.ChargingOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	order.Id = id
	err = h.chargingService.UpdateChargingOrder(ctx, &order)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})
}