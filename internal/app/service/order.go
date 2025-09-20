package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (s *ChargingService) AddChargingOrderToApplication(ctx context.Context, tariffId uint, applicationId uint) (*domain.ChargingOrder, error) {
	_, err := s.chargingRepository.GetTariffById(ctx, tariffId)
	if err != nil {
		return nil, fmt.Errorf("tariff not found: %v", err)
	}

	existingOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, applicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing orders: %v", err)
	}

	for _, order := range *existingOrders {
		if !order.IsDeleted && order.TariffId == tariffId {
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
		IsDeleted:       false,
	}

	err = s.chargingRepository.CreateChargingOrder(ctx, newChargingOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to create charging order: %v", err)
	}

	log.Printf("charging order added: application=%d, tariff=%d", applicationId, tariffId)
	return newChargingOrder, nil
}

func (s *ChargingService) RemoveChargingOrderFromApplication(ctx context.Context, orderId uint) (*domain.ChargingApplication, error) {
	existingOrder, err := s.chargingRepository.GetChargingOrderById(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("charging order not found: %v", err)
	}

	err = s.chargingRepository.DeleteChargingOrder(ctx, orderId, existingOrder.ApplicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete charging order: %v", err)
	}

	existingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, existingOrder.ApplicationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	// надо удалять черновик, если из удалили все услуги?
	if existingApplication.AmountOfOrders == 0 {
		err := s.chargingRepository.DeleteChargingApplicationById(ctx, existingOrder.ApplicationId)
		if err != nil {
			return nil, fmt.Errorf("failed to delete empty application: %v", err)
		}
	}

	log.Printf("charging order removed: id=%d, application=%d", orderId, existingOrder.ApplicationId)
	return &existingOrder.Application, nil
}

func (s *ChargingService) UpdateChargingOrder(ctx context.Context, chargingOrder *domain.ChargingOrder) error {
	_, err := s.chargingRepository.GetChargingOrderById(ctx, chargingOrder.Id)
	if err != nil {
		return fmt.Errorf("charging order not found: %v", err)
	}

	if err := s.validateChargingOrder(chargingOrder); err != nil {
		return err
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
		return fmt.Errorf("failed to update charging order: %v", err)
	}

	log.Printf("charging order updated: id=%d", chargingOrder.Id)
	return nil
}

func (s *ChargingService) validateChargingOrder(chargingOrder *domain.ChargingOrder) error { // в хелперы?
	if chargingOrder.BatteryCapacity <= 0 {
		return fmt.Errorf("battery capacity is required")
	}
	if chargingOrder.CurrentPercent < 0 {
		return fmt.Errorf("current percent is required")
	}
	return nil
}
