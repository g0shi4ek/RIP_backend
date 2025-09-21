package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (r *ChargingRepository) CreateTariff(ctx context.Context, tariff *domain.ChargingTariff) error {
	err := r.db.WithContext(ctx).Create(tariff).Error
	if err != nil {
		return fmt.Errorf("failed to create tariff: %v", err)
	}
	log.Printf("repo: tariff created: %d", tariff.Id)
	return nil
}

func (r *ChargingRepository) UpdateTariff(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error {
	if len(chargingUpdates) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Model(&domain.ChargingTariff{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(chargingUpdates).Error

	if err != nil {
		return fmt.Errorf("failed to update tariff: %v", err)
	}
	log.Printf("repo: tariff updated: %d", id)
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

	log.Println("repo: tariffList founded")
	return &tariffsList, nil
}

func (r *ChargingRepository) GetTariffById(ctx context.Context, id uint) (*domain.ChargingTariff, error) {
	var tariff domain.ChargingTariff

	err := r.db.WithContext(ctx).
		Model(&tariff).
		Where("id = ? AND is_deleted = ?", id, false).
		First(&tariff).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get tariff: %v", err)
	}
	log.Printf("repo: tariff retrieved: %d", id)
	return &tariff, nil
}

func (r *ChargingRepository) UpdateTariffImage(ctx context.Context, id uint, imageData []byte) error {
	filename := "tariff_" + strconv.Itoa(int(id)) + "image"
	imageUrl, err := r.mc.UploadTariffImage(ctx,imageData, filename)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&domain.ChargingTariff{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("image_url", imageUrl).Error

	if err != nil {
		return fmt.Errorf("failed to update tariff: %v", err)
	}
	log.Printf("repo: tariff image updated: %d", id)
	return nil
}

func (r *ChargingRepository) DeleteTariff(ctx context.Context, id uint, filename string) error {
	err := r.mc.DeleteTariffImage(ctx, filename)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).
		Model(&domain.ChargingTariff{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true).Error

	if err != nil {
		return fmt.Errorf("failed to delete tariff: %v", err)
	}
	log.Printf("repo: tariff deleted: %d", id)
	return nil
}
