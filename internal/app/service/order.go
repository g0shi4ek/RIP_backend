package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
)

func (s *ChargingService) AddChargingOrderToApplication(ctx context.Context, tariffId uint, applicationId uint) (*domain.ChargingOrder, error) {
	_, err := s.chargingRepository.GetTariffById(ctx, tariffId)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	chargingOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, applicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing orders: %v", err)
	}

	for _, chargingOrder := range *chargingOrders {
		if chargingOrder.TariffId == tariffId {
			return nil, fmt.Errorf("tariff already exists in application")
		}
	}

	newChargingOrder := &domain.ChargingOrder{
		ApplicationId:   applicationId,
		TariffId:        tariffId,
		BatteryCapacity: 0,
		CurrentPercent:  0,
		StartTime:       time.Time{},
		EstimatedTime:   0,
		CalculatedPrice: 0,
	}

	err = s.chargingRepository.CreateChargingOrder(ctx, newChargingOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to create charging order: %v", err)
	}

	existingOrder, err := s.chargingRepository.GetChargingOrderById(ctx, newChargingOrder.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	log.Printf("charging order added: application=%d, tariff=%d", applicationId, tariffId)
	return existingOrder, nil
}

func (s *ChargingService) RemoveChargingOrderFromApplication(ctx context.Context, orderId uint) (*domain.ChargingApplication, error) {
	chargingOrder, err := s.chargingRepository.GetChargingOrderById(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("charging order not found: %v", err)
	}

	err = s.chargingRepository.DeleteChargingOrder(ctx, orderId, chargingOrder.ApplicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete charging order: %v", err)
	}

	existingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, chargingOrder.ApplicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	// надо удалять черновик, если из удалили все услуги?
	/*if existingApplication.AmountOfOrders == 0 {
		err := s.chargingRepository.DeleteChargingApplicationById(ctx, chargingOrder.ApplicationId)
		if err != nil {
			return nil, fmt.Errorf("failed to delete empty application: %v", err)
		}
	}*/

	log.Printf("charging order removed: id=%d, application=%d", orderId, chargingOrder.ApplicationId)
	return existingApplication, nil
}

func (s *ChargingService) UpdateChargingOrder(ctx context.Context, chargingOrder *domain.ChargingOrder) (*domain.ChargingOrder, error) {
	_, err := s.chargingRepository.GetChargingOrderById(ctx, chargingOrder.Id)
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

	err = s.chargingRepository.UpdateChargingOrder(ctx, chargingOrder.Id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to update charging order: %v", err)
	}

	newOrder, err := s.chargingRepository.GetChargingOrderById(ctx, chargingOrder.Id)
	if err != nil {
		return nil, fmt.Errorf("charging order not found: %v", err)
	}

	log.Printf("charging order updated: id=%d", chargingOrder.Id)
	return newOrder, nil
}
