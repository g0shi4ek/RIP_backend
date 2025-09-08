package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, applicationId uint, creatorId uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication
	status := "черновик"

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("id = ? AND is_deleted = ? AND creator_id = ? AND status = ?", applicationId, false, creatorId, status).
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application with id %d not found", applicationId)
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &application, nil
}

func (r *ChargingRepository) DeleteChargingApplicationById(ctx context.Context, id uint) error {
	// SQL UPDATE без ORM
	query := "UPDATE charging_applications SET status = 'удалён', is_deleted = true, updated_at = NOW() WHERE id = ?"
	err := r.db.WithContext(ctx).Exec(query, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("application with id %d not found", id)
		}
		return fmt.Errorf("failed to delete application: %v", err)
	}
	return nil
}

// Нужна ли вообще?
func (r *ChargingRepository) UpdateChargingApplicationById(ctx context.Context, id uint, app *domain.ChargingApplication) error {
	var application domain.ChargingApplication
	creatorId := 3
	status := "черновик"

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("id = ? AND is_deleted = ? AND creator_id = ? AND status = ?", id, false, creatorId, status).
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("application with id %d not found", id)
		}
		return fmt.Errorf("failed to get application: %v", err)
	}

	updateData := map[string]interface{}{
		"status":           app.Status,
		"moderator_id":     app.ModeratorID,
		"amount_of_orders": app.AmountOfOrders,
		"total_price":      app.TotalPrice,
		"updated_at":       time.Now(),
		"completed_at":     app.CompletedAt,
	}

	err = r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Where("id = ? AND is_deleted = ? AND creator_id = ? AND status = ?", id, false, creatorId, status).
		Updates(updateData).Error

	if err != nil {
		return fmt.Errorf("failed to update application: %v", err)
	}

	return nil
}
