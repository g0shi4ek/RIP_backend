package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
)

func (s *ChargingService) GetChargingApplications(ctx context.Context, creatorId uint, userRole, status, startDate, endDate string) (*[]domain.ChargingApplication, error) {
	var startTime, endTime time.Time
	var err error

	if startDate != "" {
		startTime, err = time.Parse("02.01.2006", startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	}

	if endDate != "" {
		endTime, err = time.Parse("02.01.2006", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())
	}

	chargingApplications, err := s.chargingRepository.GetAllChargingApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging applications: %v", err)
	}

	var filteredApplications []domain.ChargingApplication
	for _, chargingApplication := range *chargingApplications {
		if (status != "" && chargingApplication.Status != status) || chargingApplication.Status == "deleted" || chargingApplication.Status == "canceled" {
			continue
		}

		if !startTime.IsZero() && chargingApplication.CreatedAt.Before(startTime) {
			continue
		}

		if !endTime.IsZero() && chargingApplication.CreatedAt.After(endTime) {
			continue
		}
		if userRole == "client" && chargingApplication.CreatorId != creatorId{
			continue
		}

		filteredApplications = append(filteredApplications, chargingApplication)
	}

	log.Printf("found %d charging applications with filters", len(filteredApplications))
	return &filteredApplications, nil
}

func (s *ChargingService) GetChargingApplication(ctx context.Context, id uint) (*domain.ChargingApplication, *[]domain.ChargingOrder, error) {
	chargingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get charging application: %v", err)
	}
	chargingApplicationOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get charging application orders: %v", err)
	}

	log.Printf("get charging application: %d", id)
	return chargingApplication, chargingApplicationOrders, nil
}

func (s *ChargingService) GetDraftChargingApplicationIfExist(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	chargingDraft, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return nil, err
	}
	return chargingDraft, nil
}

func (s *ChargingService) GetDraftChargingApplication(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	chargingDraft, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err == nil {
		return chargingDraft, nil
	}

	newChargingDraft := &domain.ChargingApplication{
		CreatorId:      creatorId,
		Status:         "draft",
		AmountOfOrders: 0,
		CreatedAt:      time.Now(),
	}

	err = s.chargingRepository.CreateChargingApplication(ctx, newChargingDraft)
	if err != nil {
		return nil, fmt.Errorf("failed to create draft: %v", err)
	}

	log.Printf("created new draft for: %d", creatorId)
	return newChargingDraft, nil
}

func (s *ChargingService) UpdateChargingApplicationPhone(ctx context.Context, phone string, creatorId uint) (*domain.ChargingApplication, error) {
	chargingApplication, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	if err := helpers.ValidatePhone(phone); err != nil {
		return nil, err
	}

	chargingUpdates := map[string]interface{}{
		"creator_phone": phone,
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, chargingApplication.Id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to update phone: %v", err)
	}
	chargingApplication.CreatorPhone = phone

	log.Printf("updated draft for: %d", creatorId)
	return chargingApplication, nil
}

func (s *ChargingService) FormChargingApplication(ctx context.Context, creatorId uint) (*domain.ChargingApplication, error) {
	chargingApplication, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application: %v", err)
	}

	if chargingApplication.CreatorPhone == "" {
		return nil, fmt.Errorf("phone number is required to form charging application")
	}

	chargingOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, chargingApplication.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging application orders: %v", err)
	}

	for _, chargingOrder := range *chargingOrders {
		if err = helpers.ValidateChargingOrder(&chargingOrder); err != nil {
			return nil, fmt.Errorf("failed to form application: %v", err)
		}
	}

	time := time.Now()
	chargingUpdates := map[string]interface{}{
		"formed_at": time,
		"status":    "formed",
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, chargingApplication.Id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to form application: %v", err)
	}
	chargingApplication.FormedAt = time
	chargingApplication.Status = "formed"

	log.Printf("formed application for: %d", creatorId)
	return chargingApplication, nil
}

func (s *ChargingService) CompleteChargingApplication(ctx context.Context, id uint, moderatorId uint) (*domain.ChargingApplication, error) {
	chargingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %v", err)
	}

	if chargingApplication.Status != "formed" {
		return nil, fmt.Errorf("can only complete formed applications")
	}

	var totalPrice float32

	chargingOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, chargingApplication.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for application: %v", err)
	}

	for _, chargingOrder := range *chargingOrders {
		tariff, err := s.chargingRepository.GetTariffById(ctx, chargingOrder.TariffId)
		if err != nil {
			return nil, fmt.Errorf("failed to get tariff for order (application=%d, tariff=%d): %v", 
				chargingOrder.ApplicationId, chargingOrder.TariffId, err)
		}

		orderCost, chargingTime, err := helpers.CalculateChargingPriceForOrder(&chargingOrder, tariff)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate price for order (application=%d, tariff=%d): %v", 
				chargingOrder.ApplicationId, chargingOrder.TariffId, err)
		}

		orderUpdates := map[string]interface{}{
			"estimated_time":   chargingTime,
			"calculated_price": orderCost,
		}

		err = s.chargingRepository.UpdateChargingOrder(ctx, chargingOrder.ApplicationId, chargingOrder.TariffId, orderUpdates)
		if err != nil {
			return nil, fmt.Errorf("failed to update order (application=%d, tariff=%d)", 
				chargingOrder.ApplicationId, chargingOrder.TariffId)
		}

		totalPrice += orderCost
	}

	chargingUpdates := map[string]interface{}{
		"status":       "completed",
		"moderator_id": moderatorId,
		"completed_at": time.Now(),
		"total_price":  totalPrice,
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to complete application: %v", err)
	}

	respApplication, _ := s.chargingRepository.GetChargingApplicationById(ctx, id)

	log.Printf("completed application %d by %d", id, moderatorId)
	return respApplication, nil
}

func (s *ChargingService) RejectChargingApplication(ctx context.Context, id uint, moderatorId uint) (*domain.ChargingApplication, error) {
	chargingApplication, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %v", err)
	}

	if chargingApplication.Status != "formed" {
		return nil, fmt.Errorf("can only complete formed applications")
	}

	chargingUpdates := map[string]interface{}{
		"status":       "rejected",
		"moderator_id": moderatorId,
	}
	chargingApplication.Status = "rejected"
	chargingApplication.ModeratorId = moderatorId

	err = s.chargingRepository.UpdateChargingApplication(ctx, id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to reject application: %v", err)
	}

	log.Printf("rejected application %d by %d", id, moderatorId)
	return chargingApplication, nil
}

func (s *ChargingService) DeleteChargingApplication(ctx context.Context, creatorId uint) error {
	chargingApplication, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return fmt.Errorf("failed to get charging application: %v", err)
	}

	err = s.chargingRepository.DeleteChargingApplicationById(ctx, chargingApplication.Id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %v", err)
	}

	log.Printf("application deleted: %d", chargingApplication.Id)
	return nil
}
