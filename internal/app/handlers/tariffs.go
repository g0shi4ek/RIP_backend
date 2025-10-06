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

// GetTarrifs godoc
// @Summary Get tariffs
// @Description Get list of charging tariffs with optional filtering
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param tariffName query string false "Filter by tariff name"
// @Success 200 {array} domain.TariffResponse
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs [get]
func (h *ChargingHandler) GetTarrifs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffs, err := h.chargingService.GetTariffs(ctx, c.Query("tariffName"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var tariffResponse []domain.TariffResponse
	for _, tariff := range *tariffs {
		tariffResponse = append(tariffResponse, tariff.ToResponse())
	}

	c.JSON(http.StatusOK, tariffResponse)
}

// GetChargingTarrifById godoc
// @Summary Get tariff by ID
// @Description Get specific charging tariff by ID
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param id path int true "Tariff ID"
// @Success 200 {object} domain.TariffResponse
// @Failure 400 {object} object "Bad request"
// @Failure 404 {object} object "Tariff not found"
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs/{id} [get]
func (h *ChargingHandler) GetChargingTarrifById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	tariff, err := h.chargingService.GetTariff(ctx, id)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, tariff.ToResponse())
}

// PostChargingTariff godoc
// @Summary Create tariff
// @Description Create new charging tariff (moderator only)
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param request body domain.TariffRequest true "Tariff data"
// @Security BearerAuth
// @Success 201 {object} domain.TariffResponse
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs [post]
func (h *ChargingHandler) PostChargingTariff(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var request domain.TariffRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	tariff := domain.ChargingTariff{
		NameofTariff: request.NameofTariff,
		Description:  request.Description,
		PricePerHour: request.PricePerHour,
		Power:        request.Power,
	}

	createdTariff, err := h.chargingService.CreateTariff(ctx, &tariff)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdTariff.ToResponse())
}

// UpdateChargingTarrifById godoc
// @Summary Update tariff
// @Description Update existing charging tariff (moderator only)
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param id path int true "Tariff ID"
// @Param request body domain.TariffRequest true "Tariff data"
// @Security BearerAuth
// @Success 200 {object} object "Tariff updated successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Tariff not found"
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs/{id} [put]
func (h *ChargingHandler) UpdateChargingTarrifById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var request domain.TariffRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	tariff := domain.ChargingTariff{
		Id:           id,
		NameofTariff: request.NameofTariff,
		Description:  request.Description,
		PricePerHour: request.PricePerHour,
		Power:        request.Power,
	}

	chargingTariff, err := h.chargingService.UpdateTariff(ctx, &tariff)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "tariff updated successfully",
		"charging_tariff": chargingTariff.ToResponse(),
	})
}

// DeleteChargingTariff godoc
// @Summary Delete tariff
// @Description Delete charging tariff (soft delete) (moderator only)
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param id path int true "Tariff ID"
// @Security BearerAuth
// @Success 200 {object} object "Tariff deleted successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Tariff not found"
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs/{id} [delete]
func (h *ChargingHandler) DeleteChargingTariff(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	err = h.chargingService.DeleteTariff(ctx, id)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tariff deleted successfully"})
}

// PostChargingTariffImage godoc
// @Summary Upload tariff image
// @Description Upload image for charging tariff (moderator only)
// @Tags Tariffs
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Tariff ID"
// @Param image formData file true "Tariff image"
// @Security BearerAuth
// @Success 200 {object} object "Image uploaded successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Tariff not found"
// @Failure 500 {object} object "Internal server error"
// @Router /tariffs/{id}/image [post]
func (h *ChargingHandler) PostChargingTariffImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		h.ErrorHandler(c, ErrImageRequired)
		return
	}
	defer file.Close()

	imageData, err := helpers.ValidateImage(header)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	chargingTariff, err := h.chargingService.UploadTariffImage(ctx, id, imageData)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "image uploaded successfully",
		"charging_tariff": chargingTariff.ToResponse(),
	})
}

// AddTariffToApplication godoc
// @Summary Add tariff to application
// @Description Add a charging tariff to current user's draft application
// @Tags Tariffs
// @Accept json
// @Produce json
// @Param id path int true "Tariff ID"
// @Security BearerAuth
// @Success 201 {object} object "Tariff added successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 403 {object} object "Forbidden"
// @Failure 404 {object} object "Tariff or application not found"
// @Failure 500 {object} object "Internal server error"
// @Router /chargingApplications/tariffs/{id} [post]
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

	_, err = h.chargingService.AddChargingOrderToApplication(ctx, tariffId, draftApplication.Id)
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to add tariff to application: %v", err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "tariff successfully added to application",
	})
}
