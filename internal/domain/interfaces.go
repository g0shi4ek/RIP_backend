package domain

import (
	"context"

	"github.com/gin-gonic/gin"
)

type IChargingRepository interface {
	GetAllTariffs(ctx context.Context) (*[]ChargingTariff, error)
	GetTariffById(ctx context.Context, id uint) (*ChargingTariff, error)
	SearchTariffs(ctx context.Context, query string) (*[]ChargingTariff, error)

	//GetChargingApplicationById(ctx context.Context, applicationId uint, creatorId uint) (*ChargingApplication, error)
	GetChargingApplicationByStatus(ctx context.Context, creatorId uint, status string) (*ChargingApplication, error)
	DeleteChargingApplicationById(ctx context.Context, applicationId uint) error
	CreateDraftChargingApplication(ctx context.Context, creatorId uint) (*ChargingApplication, error)

	CreateChargingOrder(ctx context.Context, tariffId uint, applicationId uint) error
	GetChargingOrdersByApplicationId(ctx context.Context, applicationId uint) (*[]ChargingOrder, error)
}

type IChargingHandler interface {
	RegisterChargingHandler(r *gin.Engine)
	RegisterChargingStatic(r *gin.Engine)
	ErrorChargingHandler(ctx *gin.Context, errorStatusCode int, err error)

	GetTarrifs(c *gin.Context)
	GetChargingTarrifById(c *gin.Context)

	GetChargingApplicationByStatus(c *gin.Context)
	DeleteChargingApplicationById(c *gin.Context)
	AddTariffToApplication(c *gin.Context)
}
