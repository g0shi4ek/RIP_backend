package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

/*
	CreateDraftChargingApplication(ctx context.Context, application *ChargingApplication) error
	GetChargingApplicationsWithFilters(ctx context.Context, status string, startDate, endDate time.Time) (*[]ChargingApplication, error)
	GetChargingApplicationById(ctx context.Context, id uint) (*ChargingApplication, error)
	GetDraftChargingApplicationByCreator(ctx context.Context, creatorId uint) (*ChargingApplication, error)
	UpdateChargingApplicationPhone(ctx context.Context, id uint, phone string) error
	UpdateChargingApplicationStatus(ctx context.Context, id uint, status string, moderatorId uint, completedAt time.Time) error
	UpdateChargingApplicationResult(ctx context.Context, id uint, amount float32) error
	DeleteChargingApplicationById(ctx context.Context, id uint) error
*/

func (r *ChargingRepository) CreateDraftChargingApplication(ctx context.Context, chargingApplication *domain.ChargingApplication) error {
	chargingApplication.Status = "draft"
	
	err := r.db.WithContext(ctx).Create(chargingApplication).Error
	if err != nil {
		return fmt.Errorf("failed to create draft charging application: %v", err)
	}

	return nil
}

func (r *ChargingRepository) GetDraftChargingApplicationByCreator(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("creator_id = ? AND status = ?", creatorId, "draft").
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application not found")
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &application, nil
}

func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, applicationId uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("id = ? AND NOT(status = ?)", applicationId, "deleted").
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("application not found")
		}
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	return &application, nil
}

func (r *ChargingRepository) GetChargingApplicationsWithFilters(ctx context.Context, status string, startDate, endDate time.Time) (*[]domain.ChargingApplication, error) {
	var chargingApplications []domain.ChargingApplication
	
	queryChargingApplication := r.db.WithContext(ctx).Model(&domain.ChargingApplication{}).Where("NOT(status = ?) AND NOT(status = ?)", "draft", "deleted")
	
	if status != "" {
		queryChargingApplication = queryChargingApplication.Where("status = ?", status)
	}
	
	if !startDate.IsZero() {
		queryChargingApplication = queryChargingApplication.Where("created_at >= ?", startDate)
	}
	
	if !endDate.IsZero() {
		queryChargingApplication = queryChargingApplication.Where("created_at <= ?", endDate)
	}
	
	err := queryChargingApplication.Find(&chargingApplications).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get charging applications with filters: %v", err)
	}
	
	return &chargingApplications, nil
}

func (r *ChargingRepository) UpdateChargingApplicationPhone(ctx context.Context, id uint, phone string) error {
	err := r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Where("id = ?", id).
		Update("creator_phone", phone).Error

	if err != nil {
		return fmt.Errorf("failed to update charging application status: %v", err)
	}	

	return nil
}

func (r *ChargingRepository) UpdateChargingApplicationStatus(ctx context.Context, id uint, status string, moderatorId uint, completedAt time.Time) error {
	updatedChargingApplication := map[string]interface{}{
		"status":       status,
		"moderator_id": moderatorId,
		"updated_at":   time.Now(),
	}
	
	if !completedAt.IsZero() {
		updatedChargingApplication["completed_at"] = completedAt
	}

	err := r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Where("id = ?", id).
		Updates(updatedChargingApplication).Error

	if err != nil {
		return fmt.Errorf("failed to update charging application status: %v", err)
	}

	return nil
}

func (r *ChargingRepository) UpdateChargingApplicationResult(ctx context.Context, id uint, totalChargingPrice float32) error {
	err := r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Where("id = ?", id).
		Update("total_price", totalChargingPrice).Error

	if err != nil {
		return fmt.Errorf("failed to update charging application result: %v", err)
	}

	return nil
}

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
