package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidID          = fmt.Errorf("invalid ID")
	ErrInvalidRequestBody = fmt.Errorf("invalid request body")
	ErrInvalidPhone       = fmt.Errorf("invalid phone number")
	ErrImageRequired      = fmt.Errorf("image file is required")
	ErrInvalidTariffData  = fmt.Errorf("invalid tariff data")
	ErrInvalidLogin       = fmt.Errorf("invalid user login")
	ErrInvalidPassword    = fmt.Errorf("invalid user password")
	ErrInvalidOrderData   = fmt.Errorf("invalid battery capaciry or current percent")
	ErrImageTooLarge      = fmt.Errorf("image size too large, maximum 5MB")
	ErrInvalidImageType   = fmt.Errorf("only image files are allowed")
	ErrFailedReadImage    = fmt.Errorf("failed to read image file")
	ErrInvalidDate        = fmt.Errorf("invalid date format")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
	ErrForbidden          = fmt.Errorf("forbidden")
)

type ChargingHandler struct {
	chargingService domain.IChargingService
}

func NewChargingHandler(serv domain.IChargingService) (*ChargingHandler, error) {
	return &ChargingHandler{
		chargingService: serv,
	}, nil
}

func (h *ChargingHandler) RegisterChargingHandler(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Tariffs
		api.GET("/tariffs", h.GetTarrifs)
		api.GET("/tariffs/:id", h.GetChargingTarrifById)
		api.POST("/tariffs", h.PostChargingTariff)
		api.PUT("/tariffs/:id", h.UpdateChargingTarrifById)
		api.DELETE("/tariffs/:id", h.DeleteChargingTariff)
		api.POST("/tariffs/:id/image", h.PostChargingTariffImage)

		// Applications
		api.GET("/chargingApplications", h.GetChargingApplications)
		api.GET("/chargingApplications/:id", h.GetChargingApplicationById)
		api.GET("/chargingApplications/draft", h.GetDraftChargingApplication)
		api.PUT("/chargingApplications/:id/phone", h.UpdateChargingApplicationPhone)
		api.PUT("/chargingApplications/:id/form", h.UpdateChargingApplicationByCreator)
		api.PUT("/chargingApplications/:id/:action", h.UpdateChargingApplicationByModerator) // complete/reject
		api.DELETE("/chargingApplications/:id", h.DeleteChargingApplicationById)

		// Orders
		api.POST("/chargingApplications/tariffs/:id", h.AddTariffToApplication)
		api.DELETE("/chargingOrders/:id", h.DeleteChargingOrderFromApplication)
		api.PUT("/chargingOrders/:id", h.UpdateChargingOrder)

		// Users
		api.POST("/users/register", h.RegisterUser)
		api.GET("/users/:id", h.GetUserById)
		api.PUT("/users/:id", h.UpdateUser)
		api.POST("/users/login", h.LoginUser)
		api.POST("/users/logout", h.LogOutUser)
	}
}

func (h *ChargingHandler) RegisterChargingStatic(r *gin.Engine) {}

func (h *ChargingHandler) ErrorHandler(c *gin.Context, err error) {
	logrus.Error(err.Error())

	var statusCode int

	switch {
	case errors.Is(err, ErrInvalidID) ||
		errors.Is(err, ErrInvalidRequestBody) ||
		errors.Is(err, ErrInvalidPhone) ||
		errors.Is(err, ErrInvalidTariffData) ||
		errors.Is(err, ErrInvalidOrderData) ||
		errors.Is(err, ErrImageRequired) ||
		errors.Is(err, ErrImageTooLarge) ||
		errors.Is(err, ErrInvalidImageType) ||
		errors.Is(err, ErrFailedReadImage) ||
		errors.Is(err, ErrInvalidLogin) ||
		errors.Is(err, ErrInvalidPassword) ||
		errors.Is(err, ErrInvalidDate):
		statusCode = http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(err.Error(), "already exists") ||
		strings.Contains(err.Error(), "already used"):
		statusCode = http.StatusConflict
	default:
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *ChargingHandler) getCurrentUserID(c *gin.Context) (uint, error) {
	return 3, nil
}
