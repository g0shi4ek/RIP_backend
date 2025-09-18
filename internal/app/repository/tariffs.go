package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

/*
    CreateTariff(ctx context.Context, tariff *ChargingTariff) error
    UpdateTariff(ctx context.Context, tariff *ChargingTariff) error
    UpdateTariffImage(ctx context.Context, id uint, imageUrl string) error
	GetAllTariffs(ctx context.Context) ([]ChargingTariff, error)
	GetTariffById(ctx context.Context, id uint) (*ChargingTariff, error)
	SearchTariffs(ctx context.Context, query string) ([]ChargingTariff, error)
	DeleteTariff(ctx context.Context, id uint) error
*/

func (r *ChargingRepository) CreateTariff(ctx context.Context, tariff *domain.ChargingTariff) error {
	err := r.db.WithContext(ctx).Create(tariff).Error
	if err != nil {
		return fmt.Errorf("failed to create tariff: %v", err)
	}
	log.Printf("tariff created: %d", tariff.Id)
	return nil
}

func (r *ChargingRepository) UpdateTariff(ctx context.Context, tariff *domain.ChargingTariff) error {
	err := r.db.WithContext(ctx).Model(&domain.ChargingTariff{}).
		Where("id = ? AND is_deleted = ?", tariff.Id, false).
		Updates(map[string]interface{}{
			"nameof_tariff":  tariff.NameofTariff,
			"description":    tariff.Description,
			"price_per_hour": tariff.PricePerHour,
			"power":          tariff.Power,
		}).Error

	if err != nil {
		return fmt.Errorf("failed to update tariff: %v", err)
	}

	log.Printf("tariff updated: %d", tariff.Id)
	return nil
}

func (r *ChargingRepository) GetAllTariffs(ctx context.Context) (*[]domain.ChargingTariff, error) {
	var tariffsList []domain.ChargingTariff

	err := r.db.WithContext(ctx).
		Model(&tariffsList).
		Where("is_deleted = ?", false).
		Find(&tariffsList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tariffs: %v", err)
	}

	if len(tariffsList) == 0 {
		return nil, fmt.Errorf("there are no tariffs")
	}

	log.Println("tariffList founded")
	return &tariffsList, nil
}

func (r *ChargingRepository) GetTariffById(ctx context.Context, id uint) (*domain.ChargingTariff, error) {
	var tariff domain.ChargingTariff

	err := r.db.WithContext(ctx).
		Model(&tariff).
		Where("id = ? AND is_deleted = ?", id, false).
		First(&tariff).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("tariff with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get tariff: %v", err)
	}

	return &tariff, nil
}

func (r *ChargingRepository) DeleteTariff(ctx context.Context, id uint) error {
	tariff, err := r.GetTariffById(ctx, id)
	if err != nil {
		return err
	}

	if tariff.ImageUrl != "" {
		// удаление из минио
		log.Printf("Image %s should be deleted from storage", tariff.ImageUrl)
	}

	err = r.db.WithContext(ctx).Model(&domain.ChargingTariff{}).
		Where("id = ?", id).
		Update("is_deleted", true).Error

	if err != nil {
		return fmt.Errorf("failed to delete tariff: %v", err)
	}

	log.Printf("Tariff soft deleted: %d", id)
	return nil
}

func (r *ChargingRepository) UpdateTariffImage(ctx context.Context, id uint, imageUrl string) error {
	err := r.db.WithContext(ctx).Model(&domain.ChargingTariff{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("image_url", imageUrl).Error

	if err != nil {
		return fmt.Errorf("failed to update tariff image: %v", err)
	}

	log.Printf("tariff image updated: %d", id)
	return nil
}

func (r *ChargingRepository) SearchTariffs(ctx context.Context, tariffQuery string) (*[]domain.ChargingTariff, error) {
	var filteredTariffs []domain.ChargingTariff

	tariffPattern := "%" + strings.ToLower(tariffQuery) + "%"

	err := r.db.WithContext(ctx).
		Model(&filteredTariffs).
		Where("is_deleted = ? AND ("+
			"LOWER(nameof_tariff) LIKE ? OR "+
			"LOWER(description) LIKE ?)",
			false, tariffPattern, tariffPattern).
		Find(&filteredTariffs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search tariff: %v", err)
	}

	return &filteredTariffs, nil
}
