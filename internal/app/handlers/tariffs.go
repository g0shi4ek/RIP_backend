package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
)

func (h *ChargingHandler) GetTarrifs(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tariffs, err := h.chargingService.GetTariffs(ctx, c.Query("tariffName"))
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to get tariff list: %v", err))
		return
	}

	c.JSON(http.StatusOK, tariffs)
}

func (h *ChargingHandler) GetChargingTarrifById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	tariff, err := h.chargingService.GetTariff(ctx, id)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, tariff)
}

func (h *ChargingHandler) PostChargingTariff(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var tariff domain.ChargingTariff
	if err := c.ShouldBindJSON(&tariff); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	createdTariff, err := h.chargingService.CreateTariff(ctx, &tariff)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdTariff)
}

func (h *ChargingHandler) UpdateChargingTarrifById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var tariff domain.ChargingTariff
	if err := c.ShouldBindJSON(&tariff); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	tariff.Id = id
	err = h.chargingService.UpdateTariff(ctx, &tariff)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tariff updated successfully"})
}

func (h *ChargingHandler) DeleteChargingTariff(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
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

func (h *ChargingHandler) PostChargingTariffImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := h.validateID(c.Param("id"))
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

	imageData, err := h.validateImage(header)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	err = h.chargingService.UploadTariffImage(ctx, id, imageData)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image uploaded successfully"})
}


func (h *ChargingHandler) validateImage(fileHeader *multipart.FileHeader) ([]byte, error) {
	if fileHeader.Size > 5<<20 {
		return nil, ErrImageTooLarge
	}
	
	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
		return nil, ErrInvalidImageType
	}
	
	file, err := fileHeader.Open()
	if err != nil {
		return nil, ErrFailedReadImage
	}
	defer file.Close()
	
	tariffImageData, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrFailedReadImage
	}
	
	return tariffImageData, nil
}