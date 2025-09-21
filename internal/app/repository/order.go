package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

func (r *ChargingRepository) CreateChargingOrder(ctx context.Context, chargingOrder *domain.ChargingOrder) error {
	//транзакция для согласованности данных
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&chargingOrder).Error
		if err != nil {
			return fmt.Errorf("failed to create charging order: %v", err)
		}

		err = tx.Model(&domain.ChargingApplication{}).
			Where("id = ?", chargingOrder.ApplicationId).
			Update("amount_of_orders", gorm.Expr("amount_of_orders + 1")).
			Error
		if err != nil {
			return fmt.Errorf("failed to update application order count: %v", err)
		}
		log.Printf("repo: added tariff %d to %d application", chargingOrder.TariffId, chargingOrder.ApplicationId)
		return nil
	})
}

func (r *ChargingRepository) UpdateChargingOrder(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error {
	err := r.db.WithContext(ctx).Model(&domain.ChargingOrder{}).
		Where("id = ?", id).
		Updates(chargingUpdates).Error

	if err != nil {
		return fmt.Errorf("failed to update charging order: %v", err)
	}

	log.Printf("repo: updated order, %d", id)
	return nil
}

func (r *ChargingRepository) GetChargingOrderById(ctx context.Context, id uint) (*domain.ChargingOrder, error) {
	var chargingOrder domain.ChargingOrder

	err := r.db.WithContext(ctx).
		Preload("Tariff").
		Preload("Application").
		Model(&chargingOrder).
		Where("id = ?", id).
		First(&chargingOrder).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get charging order: %v", err)
	}
	log.Printf("repo: order retrieved: %d", id)
	return &chargingOrder, nil
}

func (r *ChargingRepository) GetChargingOrdersByApplicationId(ctx context.Context, applicationId uint) (*[]domain.ChargingOrder, error) {
	var chargingOrders []domain.ChargingOrder

	err := r.db.WithContext(ctx).
		Preload("Tariff").
		Preload("Application").
		Where("application_id = ?", applicationId).
		Find(&chargingOrders).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get orders for application %d: %v", applicationId, err)
	}

	log.Printf("repo: retrieved orders for %d application", applicationId)
	return &chargingOrders, nil
}

func (r *ChargingRepository) DeleteChargingOrder(ctx context.Context, orderId uint, applicationId uint) error {
	//удаление заказа и обновление счетчика в заявке
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&domain.ChargingOrder{}).
			Where("id = ?", orderId).
			Delete(&domain.ChargingOrder{}).Error
		if err != nil {
			return fmt.Errorf("failed to delete charging order: %v", err)
		}

		err = tx.Model(&domain.ChargingApplication{}).
			Where("id = ?", applicationId).
			Update("amount_of_orders", gorm.Expr("amount_of_orders - 1")).
			Error
		if err != nil {
			return fmt.Errorf("failed to update application order count: %v", err)
		}
		log.Printf("repo: deleted order from %d application", applicationId)
		return nil
	})
}
