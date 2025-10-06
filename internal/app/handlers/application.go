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

// GetChargingApplications godoc
// @Summary Get charging applications
// @Description Get list of charging applications with optional filtering
// @Tags Charging applications
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (draft, formed, completed, rejected)"
// @Param start_date query string false "Start date for filtering"
// @Param end_date query string false "End date for filtering"
// @Security BearerAuth
// @Success 200 {array} domain.ChargingApplicationResponse
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications [get]
func (h *ChargingHandler) GetChargingApplications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	creatorId, err := h.getCurrentUserID(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}
	userRole, err := h.getCurrentUserRole(c)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}
	chargingApplications, err := h.chargingService.GetChargingApplications(ctx, creatorId, userRole,
		c.Query("status"),
		c.Query("start_date"),
		c.Query("end_date"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var chargingAppResponse []domain.ChargingApplicationResponse
	for _, app := range *chargingApplications {
		chargingAppResponse = append(chargingAppResponse, app.ToResponse())
	}

	c.JSON(http.StatusOK, chargingAppResponse)
}

// GetChargingApplicationById godoc
// @Summary Get charging application by ID
// @Description Get specific charging application with its orders
// @Tags Charging applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Security BearerAuth
// @Success 200 {object} object "Application with orders"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 404 {object} object "Application not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/{id} [get]
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
		"charging_application": resultApplication.ToResponse(),
		"charging_orders":      orderResponses,
	}

	c.JSON(http.StatusOK, chargingResponse)
}

// GetDraftChargingApplication godoc
// @Summary Get draft charging application
// @Description Get current user's draft charging application
// @Tags Charging applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.ChargingDraftResponse
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/draft [get]
func (h *ChargingHandler) GetDraftChargingApplication(c *gin.Context) {
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

	c.JSON(http.StatusOK, draftApplication.ToDraftResponse())
}

// UpdateChargingApplicationPhone godoc
// @Summary Update application phone number
// @Description Update phone number for current user's draft application
// @Tags Charging applications
// @Accept json
// @Produce json
// @Param request body domain.PhoneRequest true "Phone number"
// @Security BearerAuth
// @Success 200 {object} object "Phone updated successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/phone [put]
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
		"message":              "phone updated successfully",
		"charging_application": chargingApplication.ToResponse(),
	})
}

// UpdateChargingApplicationByCreator godoc
// @Summary Form charging application
// @Description Form (submit) the draft charging application
// @Tags Charging applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Application formed successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/form [put]
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
		"message":              "application formed successfully",
		"charging_application": chargingApplication.ToResponse(),
	})
}

// UpdateChargingApplicationByModerator godoc
// @Summary Update application status by moderator
// @Description Complete or reject charging application (moderator only)
// @Tags Charging applications
// @Accept json
// @Produce json
// @Param id path int true "Application ID"
// @Param action path string true "Action (complete, reject)"
// @Security BearerAuth
// @Success 200 {object} object "Application status updated"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Application not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/{id}/{action} [put]
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
		"message":              fmt.Sprintf("application %s successfully", action),
		"charging_application": chargingApplication.ToResponse(),
	})
}

// DeleteChargingApplicationById godoc
// @Summary Delete charging application
// @Description Delete current user's draft charging application
// @Tags Charging applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Application deleted successfully"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Application not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications [delete]
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
