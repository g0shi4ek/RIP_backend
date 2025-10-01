package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
)

func (s *ChargingService) CreateTariff(ctx context.Context, tariff *domain.ChargingTariff) (*domain.ChargingTariff, error) {
	if err := helpers.ValidateTariff(tariff); err != nil {
		return nil, err
	}

	err := s.chargingRepository.CreateTariff(ctx, tariff)
	if err != nil {
		return nil, fmt.Errorf("failed to create tariff: %v", err)
	}

	log.Printf("tariff created: %d", tariff.Id)
	return tariff, nil
}

func (s *ChargingService) UpdateTariff(ctx context.Context, tariff *domain.ChargingTariff) (*domain.ChargingTariff, error) {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, tariff.Id)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	if err := helpers.ValidateTariff(tariff); err != nil {
		return nil, err
	}

	chargingUpdates := map[string]interface{}{
		"nameof_tariff":  tariff.NameofTariff,
		"description":    tariff.Description,
		"price_per_hour": tariff.PricePerHour,
		"power":          tariff.Power,
	}

	err = s.chargingRepository.UpdateTariff(ctx, existingTariff.Id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to update tariff data: %v", err)
	}
	newTariff, err := s.chargingRepository.GetTariffById(ctx, tariff.Id)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	log.Printf("updated tariff: %d", tariff.Id)
	return newTariff, nil
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

	imageList := strings.Split(existingTariff.ImageUrl, "/")
	filename := imageList[len(imageList)-1]

	err = s.chargingRepository.DeleteTariff(ctx, existingTariff.Id, filename)
	if err != nil {
		return fmt.Errorf("failed to delete tariff: %v", err)
	}

	log.Printf("tariff deleted: %d", id)
	return nil
}

func (s *ChargingService) UploadTariffImage(ctx context.Context, id uint, tariffImage []byte) (*domain.ChargingTariff, error) {
	existingTariff, err := s.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	err = s.chargingRepository.UpdateTariffImage(ctx, existingTariff.Id, tariffImage)
	if err != nil {
		return nil, fmt.Errorf("failed to upload tariff image: %v", err)
	}

	newTariff, err := s.chargingRepository.GetTariffById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	log.Printf("image uploaded: %d", id)
	return newTariff, nil
}