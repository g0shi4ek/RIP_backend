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

func (h *ChargingHandler) GetChargingApplications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	chargingApplications, err := h.chargingService.GetChargingApplications(ctx, 
		c.Query("status"), 
		c.Query("start_date"), 
		c.Query("end_date"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var chargingAppResponse []domain.ChargingApplicationResponse
    for _, app := range *chargingApplications{
        chargingAppResponse = append(chargingAppResponse, app.ToResponse())
    }

    c.JSON(http.StatusOK, chargingAppResponse)
}

func (h *ChargingHandler) GetChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	resultApplication, chargingOrders, err := h.chargingService.GetChargingApplication(ctx, id)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}
	var orderResponses []domain.ChargingOrderResponse
    for _, order := range *chargingOrders {
        orderResponses = append(orderResponses, order.ToResponse())
    }

    chargingResponse := gin.H{
        "application": resultApplication.ToResponse(),
        "orders":      orderResponses,
    }

    c.JSON(http.StatusOK, chargingResponse)
}

func (h *ChargingHandler) GetDraftChargingApplication(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

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

	c.JSON(http.StatusOK, draftApplication.ToResponse())
}

func (h *ChargingHandler) UpdateChargingApplicationPhone(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var phoneRequest domain.PhoneRequest
	if err := c.ShouldBindJSON(&phoneRequest); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	chargingApplication, err := h.chargingService.UpdateChargingApplicationPhone(ctx, phoneRequest.Phone, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "phone updated successfully",
		"application": chargingApplication.ToResponse(),
	})
}

func (h *ChargingHandler) UpdateChargingApplicationByCreator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	chargingApplication, err := h.chargingService.FormChargingApplication(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to form application: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "application formed successfully",
		"application": chargingApplication.ToResponse(),
	})
}

func (h *ChargingHandler) UpdateChargingApplicationByModerator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	moderatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	action := c.Param("action")
	var chargingApplication *domain.ChargingApplication
	switch action {
	case "complete":
		chargingApplication, err = h.chargingService.CompleteChargingApplication(ctx, id, moderatorId)
	case "reject":
		chargingApplication, err = h.chargingService.RejectChargingApplication(ctx, id, moderatorId)
	default:
		h.ErrorHandler(c, fmt.Errorf("invalid action: %s", action))
		return
	}

	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("application %s successfully", action),
		"application": chargingApplication.ToResponse(),
	})
}

func (h *ChargingHandler) DeleteChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	err = h.chargingService.DeleteChargingApplication(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}