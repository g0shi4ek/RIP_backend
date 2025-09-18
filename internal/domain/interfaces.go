package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type IChargingRepository interface {
	GetAllTariffs(ctx context.Context) (*[]ChargingTariff, error)
	GetTariffById(ctx context.Context, id uint) (*ChargingTariff, error)
	SearchTariffs(ctx context.Context, query string) (*[]ChargingTariff, error) // перенести выше?
	CreateTariff(ctx context.Context, tariff *ChargingTariff) error
	UpdateTariff(ctx context.Context, tariff *ChargingTariff) error
	UpdateTariffImage(ctx context.Context, id uint, imageUrl string) error
	DeleteTariff(ctx context.Context, id uint) error

	GetAllChargingApplications(ctx context.Context) (*[]ChargingApplication, error)
	GetChargingApplicationsWithFilters(ctx context.Context, status string, startDate, endDate time.Time) (*[]ChargingApplication, error) // перенести выше?
	GetChargingApplicationById(ctx context.Context, id uint) (*ChargingApplication, error)
	GetDraftChargingApplicationByCreator(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	CreateDraftChargingApplication(ctx context.Context, application *ChargingApplication) error
	UpdateChargingApplicationPhone(ctx context.Context, application *ChargingApplication) error
	UpdateChargingApplicationStatus(ctx context.Context, id uint, status string, moderatorId uint, completedAt time.Time) error
	UpdateChargingApplicationResult(ctx context.Context, id uint, amount float32) error
	DeleteChargingApplication(ctx context.Context, id uint) error

	CreateChargingOrder(ctx context.Context, tariffId uint, applicationId uint) error
	GetChargingOrdersByApplicationId(ctx context.Context, applicationId uint) (*[]ChargingOrder, error)
	UpdateChargingOrder(ctx context.Context, order *ChargingOrder) error
	DeleteChargingOrder(ctx context.Context, id uint) error

	GetChargingUserById(ctx context.Context, id uint) (*User, error)
	GetChargingUserByLogin(ctx context.Context, login string) (*User, error)
	CreateChargingUser(ctx context.Context, user *User) error
	UpdateChargingUser(ctx context.Context, user *User) error
}

type IChargingService interface {
	GetTariffs(ctx context.Context, searchQuery string) (*[]ChargingTariff, error)
	GetTariff(ctx context.Context, id uint) (*ChargingTariff, error)
	CreateTariff(ctx context.Context, tariff *ChargingTariff) (*ChargingTariff, error)
	UpdateTariff(ctx context.Context, id uint, tariff *ChargingTariff) (*ChargingTariff, error)
	DeleteTariff(ctx context.Context, id uint) error
	UploadTariffImage(ctx context.Context, id uint, imageData []byte, fileName string) (*ChargingTariff, error)

	GetApplications(ctx context.Context, status, startDate, endDate string) (*[]ChargingApplication, error)
	GetApplication(ctx context.Context, id uint) (*ChargingApplication, error)
	GetDraftApplication(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	UpdateApplicationPhone(ctx context.Context, id uint, phone string, creatorId uint) (*ChargingApplication, error)
	FormApplication(ctx context.Context, id uint, creatorId uint) (*ChargingApplication, error)
	CompleteApplication(ctx context.Context, id uint, moderatorId uint) (*ChargingApplication, error)
	RejectApplication(ctx context.Context, id uint, moderatorId uint) (*ChargingApplication, error)
	DeleteApplication(ctx context.Context, id uint, creatorId uint) error

	AddTariffToApplication(ctx context.Context, request *AddTariffRequest, creatorId uint) (*ChargingOrder, error)
	RemoveOrderFromApplication(ctx context.Context, orderId uint, creatorId uint) error
	UpdateOrder(ctx context.Context, orderId uint, request *UpdateOrderRequest, creatorId uint) (*ChargingOrder, error)
	CalculateOrderPrice(order *ChargingOrder, tariff *ChargingTariff) error // мб в хелперы

	RegisterUser(ctx context.Context, user *User) (*User, error)
	GetUserProfile(ctx context.Context, userId uint) (*User, error)
	UpdateUserProfile(ctx context.Context, userId uint, userData *User) (*User, error)
	AuthenticateUser(ctx context.Context, login, password string) (*User, string, error)
	LogoutUser(ctx context.Context, token string) error
	ValidateUserPermissions(ctx context.Context, userId uint, applicationId uint) (bool, error)
}

type IChargingHandler interface {
	RegisterChargingHandler(r *gin.Engine)
	RegisterChargingStatic(r *gin.Engine)
	ErrorChargingHandler(ctx *gin.Context, errorStatusCode int, err error)

	GetTarrifs(c *gin.Context)               // GET /api/tariffs?tariffName=
	GetChargingTarrifById(c *gin.Context)    // GET /api/tariffs/:id
	PostChargingTariff(c *gin.Context)       // POST /api/tariffs
	UpdateChargingTarrifById(c *gin.Context) // PUT /api/tariffs/:id
	DeleteChargingTariff(c *gin.Context)     // DELETE /api/tariffs/:id
	PostChargingTariffImage(c *gin.Context)  // POST /api/tariffs/:id/image

	GetChargingApplications(c *gin.Context)          // GET /api/applications?status=&start_date=&end_date=
	GetChargingApplicationById(c *gin.Context)       // GET /api/applications/:id
	GetChargingApplicationByStatus(c *gin.Context)   // GET /api/applications/draft (корзина)
	UpdateChargingApplicationPhone(c *gin.Context)   // PUT /api/applications/:id/phone
	UpdateChargingApplicationByCreator(c *gin.Context) // PUT /api/applications/:id/form (сформировать)
	UpdateChargingApplicationByModerator(c *gin.Context) // PUT /api/applications/:id/complete или /reject
	DeleteChargingApplicationById(c *gin.Context)    // DELETE /api/applications/:id

	AddTariffToApplication(c *gin.Context)           // POST /api/applications/:applicationId/orders
	DeleteChargingOrderFromApplication(c *gin.Context) // DELETE /api/orders/:id
	UpdateChargingOrder(c *gin.Context)              // PUT /api/orders/:id

	RegisterUser(c *gin.Context)                     // POST /api/auth/register
	GetUserById(c *gin.Context)                      // GET /api/user/profile
	UpdateUser(c *gin.Context)                       // PUT /api/user/profile
	LoginUser(c *gin.Context)                        // POST /api/auth/login
	LogOutUser(c *gin.Context)                       // POST /api/auth/logout

	GetChargingApplicationCart(c *gin.Context)       // GET /api/cart (иконка корзины)
}