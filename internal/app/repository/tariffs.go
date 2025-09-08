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
