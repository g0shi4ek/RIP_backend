package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

func (r *ChargingRepository) CreateChargingOrder(ctx context.Context, tariffId uint, applicationId uint) error {
	// знаем, что существует тариф и заявка (слой бизнес логики в хендлер?)
	var existingOrder domain.ChargingOrder
	err := r.db.WithContext(ctx).
		Where("application_id = ? AND tariff_id = ? AND is_deleted = ?", applicationId, tariffId, false).
		First(&existingOrder).Error

	if err == nil {
		return fmt.Errorf("tariff already exists in application")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing order: %v", err)
	}

	newChargingOrder := domain.ChargingOrder{
		ApplicationId:   applicationId,
		TariffId:        tariffId,
		BatteryCapacity: 0,
		CurrentPercent:  0,
		StartTime:       time.Time{},
		EstimatedTime:   0,
		CalculatedPrice: 0,
		IsDeleted:       false,
	}

	//транзакция для согласованности данных
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&newChargingOrder).Error
		if err != nil {
			return fmt.Errorf("failed to create charging order: %v", err)
		}

		err = tx.Model(&domain.ChargingApplication{}).
			Where("id = ?", applicationId).
			Update("amount_of_orders", gorm.Expr("amount_of_orders + 1")).
			Error
		if err != nil {
			return fmt.Errorf("failed to update application order count: %v", err)
		}

		return nil
	})
}

func (r *ChargingRepository) GetChargingOrdersByApplicationId(ctx context.Context, applicationId uint) (*[]domain.ChargingOrder, error) {
	var chargingOrders []domain.ChargingOrder

	err := r.db.WithContext(ctx).
		Preload("Tariff").
		Preload("Application").
		Where("application_id = ? AND is_deleted = ?", applicationId, false).
		Find(&chargingOrders).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get orders for application %d: %v", applicationId, err)
	}

	return &chargingOrders, nil
}
