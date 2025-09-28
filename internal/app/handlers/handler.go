package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/app/middleware"
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
	authMiddleware  *middleware.AuthMiddleware
}

func NewChargingHandler(serv domain.IChargingService) (*ChargingHandler, error) {
	authMiddleware := middleware.NewAuthMiddleware(serv)

	return &ChargingHandler{
		chargingService: serv,
		authMiddleware:  authMiddleware,
	}, nil
}

func (h *ChargingHandler) RegisterChargingHandler(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/users/register", h.RegisterUser)
		api.POST("/users/login", h.LoginUser)
		api.GET("/tariffs", h.GetTarrifs)
		api.GET("/tariffs/:id", h.GetChargingTarrifById)

		protected := api.Group("")
		protected.Use(h.authMiddleware.Auth())
		{
			protected.GET("/users/:id", h.GetUserById)
			protected.PUT("/users/:id", h.UpdateUser)
			protected.POST("/users/logout", h.LogOutUser)
			protected.GET("/chargingApplications", h.GetChargingApplications)
			protected.GET("/chargingApplications/:id", h.GetChargingApplicationById)

			client := protected.Group("")
			client.Use(h.authMiddleware.Role("client"))
			{
				client.GET("/chargingApplications/draft", h.GetDraftChargingApplication)
				client.PUT("/chargingApplications/phone", h.UpdateChargingApplicationPhone)
				client.PUT("/chargingApplications/form", h.UpdateChargingApplicationByCreator)
				client.DELETE("/chargingApplications", h.DeleteChargingApplicationById)

				client.POST("/chargingApplications/tariffs/:id", h.AddTariffToApplication)
				client.DELETE("/chargingOrders/:id", h.DeleteChargingOrderFromApplication)
				client.PUT("/chargingOrders/:id", h.UpdateChargingOrder)
			}

			moderator := protected.Group("")
			moderator.Use(h.authMiddleware.Role("moderator"))
			{
				moderator.POST("/tariffs", h.PostChargingTariff)
				moderator.PUT("/tariffs/:id", h.UpdateChargingTarrifById)
				moderator.DELETE("/tariffs/:id", h.DeleteChargingTariff)
				moderator.POST("/tariffs/:id/image", h.PostChargingTariffImage)
				moderator.PUT("/chargingApplications/:id/:action", h.UpdateChargingApplicationByModerator)
			}
		}
	}
}

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
		errors.Is(err, ErrInvalidDate) ||
		strings.Contains(err.Error(), "invalid"):
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
	return middleware.GetCurrentUserID(c)
}
