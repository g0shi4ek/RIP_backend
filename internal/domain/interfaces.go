package domain

import (
	"context"
)

type IChargingRepository interface {
	GetTariffById(ctx context.Context, id int) ( *ChargingTariff, error)
	GetChargingApplicationById(ctx context.Context, id int) (*ChargingApplication, error)
	GetAllTariffs (ctx context.Context) (*[]ChargingTariff, error)
	SearchTariffs(ctx context.Context, query string) (*[]ChargingTariff, error)
}