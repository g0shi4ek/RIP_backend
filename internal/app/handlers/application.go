package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	c.JSON(http.StatusOK, chargingApplications)
}

func (h *ChargingHandler) GetChargingApplicationById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	resultApplication, chargingOrders, err := h.chargingService.GetChargingApplication(ctx, id)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}
	chargingResponse := gin.H{
		"application": resultApplication,
		"orders":      chargingOrders,
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

	c.JSON(http.StatusOK, draftApplication)
}

func (h *ChargingHandler) UpdateChargingApplicationPhone(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var request struct {
		Phone string `json:"phone" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	err = h.chargingService.UpdateChargingApplicationPhone(ctx, request.Phone, creatorId)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "phone updated successfully"})
}

func (h *ChargingHandler) UpdateChargingApplicationByCreator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	err = h.chargingService.FormChargingApplication(ctx, creatorId)
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to form application: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application formed successfully"})
}

func (h *ChargingHandler) UpdateChargingApplicationByModerator(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
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
	switch action {
	case "complete":
		err = h.chargingService.CompleteChargingApplication(ctx, id, moderatorId)
	case "reject":
		err = h.chargingService.RejectChargingApplication(ctx, id, moderatorId)
	default:
		h.ErrorHandler(c, fmt.Errorf("invalid action: %s", action))
		return
	}

	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("application %s successfully", action)})
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