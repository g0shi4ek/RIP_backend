package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (r *ChargingRepository) CreateChargingApplication(ctx context.Context, chargingApplication *domain.ChargingApplication) error {
	err := r.db.WithContext(ctx).Create(chargingApplication).Error
	if err != nil {
		return fmt.Errorf("failed to create draft charging application: %v", err)
	}

	log.Printf("repo: application created, %d", chargingApplication.Id)
	return nil
}

func (r *ChargingRepository) UpdateChargingApplication(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error {
	if len(chargingUpdates) == 0 {
		return nil
	}

	chargingUpdates["updated_at"] = time.Now()

	err := r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Where("id = ?", id).
		Updates(chargingUpdates).Error

	if err != nil {
		return fmt.Errorf("failed to update charging application: %v", err)
	}

	log.Printf("repo: application updated, %d", id)
	return nil
}


func (r *ChargingRepository) GetDraftChargingApplicationByCreator(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("creator_id = ? AND status = ?", creatorId, "draft").
		First(&application).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	log.Printf("repo: draft application for %d retrieved", creatorId)
	return &application, nil
}

func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, id uint) (*domain.ChargingApplication, error) {
	var application domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&application).
		Where("id = ? AND NOT(status = ?)", id, "deleted").
		First(&application).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get application: %v", err)
	}

	log.Printf("repo: application retrieved, %d", id)
	return &application, nil
}

func (r *ChargingRepository) GetAllChargingApplications(ctx context.Context) (*[]domain.ChargingApplication, error) {
	var chargingApplications []domain.ChargingApplication

	err := r.db.WithContext(ctx).
		Model(&domain.ChargingApplication{}).
		Find(&chargingApplications).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get charging applications: %v", err)
	}

	log.Println("repo: get all applications")
	return &chargingApplications, nil
}


func (r *ChargingRepository) DeleteChargingApplicationById(ctx context.Context, id uint) error {
	// SQL UPDATE без ORM
	query := "UPDATE charging_applications SET status = 'deleted', updated_at = NOW() WHERE id = ?"
	err := r.db.WithContext(ctx).Exec(query, id).Error

	if err != nil {
		return fmt.Errorf("failed to delete application: %v", err)
	}

	log.Printf("repo: application deleted (soft), %d", id)
	return nil
}
