package service

import (
	"context"
	"fmt"
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
)

func (s *ChargingService) RemoveChargingOrderFromApplication(ctx context.Context, applicationId, tariffId uint) (*domain.ChargingApplication, error) {
	_, err := s.chargingRepository.GetChargingOrder(ctx, applicationId, tariffId)
	if err != nil {
		return nil, fmt.Errorf("charging order not found: %v", err)
	}

	err = s.chargingRepository.DeleteChargingOrder(ctx, applicationId, tariffId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete charging order: %v", err)
	}

	existingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, applicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	log.Printf("charging order removed: application=%d, tariff=%d", applicationId, tariffId)
	return existingApplication, nil
}

func (s *ChargingService) UpdateChargingOrder(ctx context.Context, chargingOrder *domain.ChargingOrder) (*domain.ChargingOrder, error) {
	_, err := s.chargingRepository.GetChargingOrder(ctx, chargingOrder.ApplicationId, chargingOrder.TariffId)
	if err != nil {
		return nil, fmt.Errorf("charging order not found: %v", err)
	}

	if err := helpers.ValidateChargingOrder(chargingOrder); err != nil {
		return nil, err
	}

	chargingUpdates := map[string]interface{}{
		"battery_capacity": chargingOrder.BatteryCapacity,
		"current_percent":  chargingOrder.CurrentPercent,
		"start_time":       chargingOrder.StartTime,
		"estimated_time":   chargingOrder.EstimatedTime,
		"calculated_price": chargingOrder.CalculatedPrice,
	}

	err = s.chargingRepository.UpdateChargingOrder(ctx, chargingOrder.ApplicationId, chargingOrder.TariffId, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to update charging order: %v", err)
	}

	updatedOrder, err := s.chargingRepository.GetChargingOrder(ctx, chargingOrder.ApplicationId, chargingOrder.TariffId)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated charging order: %v", err)
	}

	log.Printf("charging order updated: application=%d, tariff=%d", chargingOrder.ApplicationId, chargingOrder.TariffId)
	return updatedOrder, nil
}
