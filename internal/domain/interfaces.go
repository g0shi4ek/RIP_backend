package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type ISessionStorage interface {
	AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsInBlacklist(ctx context.Context, token string) (bool, error)
}

type ITariffImagesStorage interface {
	UploadTariffImage(ctx context.Context, imageData []byte, filename string) (string, error)
	DeleteTariffImage(ctx context.Context, filename string) error
}

type IChargingRepository interface {
	CreateTariff(ctx context.Context, tariff *ChargingTariff) error
	UpdateTariff(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error
	GetAllTariffs(ctx context.Context) (*[]ChargingTariff, error)
	GetTariffById(ctx context.Context, id uint) (*ChargingTariff, error)
	UpdateTariffImage(ctx context.Context, id uint, imageData []byte) error
	DeleteTariff(ctx context.Context, id uint, filename string) error

	CreateChargingApplication(ctx context.Context, chargingApplication *ChargingApplication) error
	UpdateChargingApplication(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error
	GetDraftChargingApplicationByCreator(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	GetChargingApplicationById(ctx context.Context, id uint) (*ChargingApplication, error)
	GetAllChargingApplications(ctx context.Context) (*[]ChargingApplication, error)
	DeleteChargingApplicationById(ctx context.Context, id uint) error

	CreateChargingOrder(ctx context.Context, chargingOrder *ChargingOrder) error
	UpdateChargingOrder(ctx context.Context, applicationId, tariffId uint, chargingUpdates map[string]interface{}) error
	GetChargingOrder(ctx context.Context, applicationId, tariffId uint) (*ChargingOrder, error)
	GetChargingOrdersByApplicationId(ctx context.Context, applicationId uint) (*[]ChargingOrder, error)
	DeleteChargingOrder(ctx context.Context, applicationId, tariffId uint) error

	CreateChargingUser(ctx context.Context, user *User) error
	UpdateChargingUser(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error
	GetChargingUserById(ctx context.Context, id uint) (*User, error)
	GetChargingUserByLogin(ctx context.Context, login string) (*User, error)

	AddTokenToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type IChargingService interface {
	GetTariffs(ctx context.Context, tariffName string) (*[]ChargingTariff, error)
	GetTariff(ctx context.Context, id uint) (*ChargingTariff, error)
	CreateTariff(ctx context.Context, tariff *ChargingTariff) (*ChargingTariff, error)
	UpdateTariff(ctx context.Context, tariff *ChargingTariff) (*ChargingTariff, error)
	DeleteTariff(ctx context.Context, id uint) error
	UploadTariffImage(ctx context.Context, id uint, tariffImage []byte) (*ChargingTariff, error)
	AddChargingOrderToApplication(ctx context.Context, tariffId uint, applicationId uint) (*ChargingOrder, error)

	GetChargingApplications(ctx context.Context, creatorId uint, userRole string, status, startDate, endDate string) (*[]ChargingApplication, error)
	GetChargingApplication(ctx context.Context, id uint) (*ChargingApplication, *[]ChargingOrder, error)
	GetDraftChargingApplication(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	GetDraftChargingApplicationIfExist(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	UpdateChargingApplicationPhone(ctx context.Context, phone string, creatorId uint) (*ChargingApplication, error)
	FormChargingApplication(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	CompleteChargingApplication(ctx context.Context, id uint, moderatorId uint) (*ChargingApplication, error)
	RejectChargingApplication(ctx context.Context, id uint, moderatorId uint) (*ChargingApplication, error)
	DeleteChargingApplication(ctx context.Context, creatorId uint) error

	RemoveChargingOrderFromApplication(ctx context.Context, applicationId, tariffId uint) (*ChargingApplication, error)
	UpdateChargingOrder(ctx context.Context, chargingOrder *ChargingOrder) (*ChargingOrder, error)
	UpdateOrderCalculation(ctx context.Context, applicationId, tariffId uint, updates map[string]interface{}) error
	UpdateApplicationTotalPrice(ctx context.Context, applicationId uint) error

	RegisterChargingUser(ctx context.Context, user *User) (*User, error)
	GetChargingUserProfile(ctx context.Context, id uint) (*User, error)
	UpdateChargingUserProfile(ctx context.Context, id uint, userData *User) (*User, error)
	LoginChargingUser(ctx context.Context, login, password string) (string, error)
	LogoutChargingUser(ctx context.Context, token string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type IChargingHandler interface {
	RegisterChargingHandler(r *gin.Engine)
	ErrorHandler(c *gin.Context, err error)

	GetTarrifs(c *gin.Context)               // GET /api/tariffs?tariffName=
	GetChargingTarrifById(c *gin.Context)    // GET /api/tariffs/:id
	PostChargingTariff(c *gin.Context)       // POST /api/tariffs
	UpdateChargingTarrifById(c *gin.Context) // PUT /api/tariffs/:id
	DeleteChargingTariff(c *gin.Context)     // DELETE /api/tariffs/:id
	PostChargingTariffImage(c *gin.Context)  // POST /api/tariffs/:id/image
	AddTariffToApplication(c *gin.Context)   // POST /api/chargingApplications/tariffs/:id

	GetChargingApplications(c *gin.Context)              // GET /api/chargingApplications?status=&start_date=&end_date=
	GetChargingApplicationById(c *gin.Context)           // GET /api/chargingApplications/:id
	GetDraftChargingApplication(c *gin.Context)          // GET /api/chargingApplications/draft (корзина)
	UpdateChargingApplicationPhone(c *gin.Context)       // PUT /api/chargingApplications/:id/phone
	UpdateChargingApplicationByCreator(c *gin.Context)   // PUT /api/chargingApplications/:id/form (сформировать)
	UpdateChargingApplicationByModerator(c *gin.Context) // PUT /api/chargingApplications/:id/complete или /reject
	DeleteChargingApplicationById(c *gin.Context)        // DELETE /api/chargingApplications/:id

	DeleteChargingOrderFromApplication(c *gin.Context) // DELETE /api/chargingOrders/:id
	UpdateChargingOrder(c *gin.Context)                // PUT /api/chargingOrders/:id
	UpdateOrderCalculation(c *gin.Context)

	RegisterUser(c *gin.Context) // POST /api/users/register
	GetUserById(c *gin.Context)  // GET /api/users/:userId
	UpdateUser(c *gin.Context)   // PUT /api/users/:userId
	LoginUser(c *gin.Context)    // POST /api/users/login
	LogOutUser(c *gin.Context)   // POST /api/users/logout
}
