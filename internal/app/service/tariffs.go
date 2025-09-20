package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (s *ChargingService) CreateTariff(ctx context.Context, tariff *domain.ChargingTariff) (*domain.ChargingTariff, error) {
	if err := s.validateTariff(tariff); err != nil {
		return nil, err
	}

	err := s.chargingRepository.CreateTariff(ctx, tariff)
	if err != nil {
		return nil, fmt.Errorf("failed to create tariff: %v", err)
	}

	log.Printf("tariff created: %d", tariff.Id)
	return tariff, nil
}

func (s *ChargingService) UpdateTariff(ctx context.Context, tariff *domain.ChargingTariff) error {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, tariff.Id)
	if err != nil {
		return fmt.Errorf("tariff not found: %v", err)
	}

	if err := s.validateTariff(tariff); err != nil {
		return err
	}

	chargingUpdates := map[string]interface{}{
		"nameof_tariff":  tariff.NameofTariff,
		"description":    tariff.Description,
		"price_per_hour": tariff.PricePerHour,
		"power":          tariff.Power,
	}

	err = s.chargingRepository.UpdateTariff(ctx, existingTariff.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to update tariff data: %v", err)
	}

	log.Printf("updated tariff: %d", tariff.Id)
	return nil
}

func (s *ChargingService) GetTariffs(ctx context.Context, tariffName string) (*[]domain.ChargingTariff, error) {
	tariffs, err := s.chargingRepository.GetAllTariffs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariffs: %v", err)
	}

	if tariffName == "" {
		return tariffs, nil
	}

	patternTariff := strings.ToLower(tariffName)
	var filteredTariffs []domain.ChargingTariff
	for _, tariff := range *tariffs {
		name := strings.ToLower(tariff.NameofTariff)
		desc := strings.ToLower(tariff.Description)
		if strings.Contains(name, patternTariff) ||  strings.Contains(desc, patternTariff){
			filteredTariffs = append(filteredTariffs, tariff)
		}
	}

	log.Printf("found %d tariffs matching name: %s", len(filteredTariffs), tariffName)
	return &filteredTariffs, nil
}

func (s *ChargingService) GetTariff(ctx context.Context, id uint) (*domain.ChargingTariff, error) {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	log.Printf("get tariff: %d", id)
	return existingTariff, nil
}


func (s *ChargingService) DeleteTariff(ctx context.Context, id uint) error {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		return fmt.Errorf("tariff not found: %v", err)
	}

	if existingTariff.ImageUrl != ""{
		// удаление из минио
	}

	chargingUpdates := map[string]interface{}{
		"is_deleted":  true,
	}

	err = s.chargingRepository.UpdateTariff(ctx, existingTariff.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to delete tariff: %v", err)
	}

	log.Printf("tariff deleted: %d", id)
	return nil
}

func (s *ChargingService) UploadTariffImage(ctx context.Context, id uint, tariffImage []byte) error {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		return fmt.Errorf("tariff not found: %v", err)
	}

	// добавление в минио => возврат урла
	tariffImageUrl := "www"

	chargingUpdates := map[string]interface{}{
		"image_url": tariffImageUrl,
	}

	err = s.chargingRepository.UpdateTariff(ctx, existingTariff.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to upload tariff image: %v", err)
	}

	log.Printf("image uploaded: %d", id)
	return nil
}

func (s *ChargingService) validateTariff(tariff *domain.ChargingTariff) error { // в хелперы?
	if tariff.NameofTariff == "" {
		return fmt.Errorf("tariff name is required")
	}
	if tariff.Description == "" {
		return fmt.Errorf("tariff description is required")
	}
	if tariff.PricePerHour <= 0 {
		return fmt.Errorf("price per hour must be positive")
	}
	if tariff.Power <= 0 {
		return fmt.Errorf("power must be positive")
	}
	return nil
}