package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

func (r *ChargingRepository) CreateDraftChargingApplication(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	newApplication := &domain.ChargingApplication{
		CreatorID:      creatorId,
		ModeratorID:    1,
		Status:         "draft",
		AmountOfOrders: 0,
		TotalPrice:     0,
	}

	err := r.db.WithContext(ctx).Create(newApplication).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create draft application: %v", err)
	}

	return newApplication, nil
}

func (r *ChargingRepository) GetChargingApplicationByStatus(ctx context.Context, creatorId uint, status string) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("creator_id = ? AND status = ?", creatorId, status).
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application with status %s not found", status)
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &application, nil
}

/*func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, applicationId uint, creatorId uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("id = ? AND NOT(status = ?) AND creator_id = ?", applicationId, "deleted", creatorId).
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application with id %d not found", applicationId)
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &application, nil
}*/

func (r *ChargingRepository) DeleteChargingApplicationById(ctx context.Context, id uint) error {
	// SQL UPDATE без ORM
	query := "UPDATE charging_applications SET status = 'deleted', updated_at = NOW() WHERE id = ?"
	err := r.db.WithContext(ctx).Exec(query, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("application with id %d not found", id)
		}
		return fmt.Errorf("failed to delete application: %v", err)
	}
	return nil
}
