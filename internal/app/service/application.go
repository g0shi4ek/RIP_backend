package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (s *ChargingService) GetChargingApplications(ctx context.Context, status, startDate, endDate string) (*[]domain.ChargingApplication, error) {
	var startTime, endTime time.Time
	var err error

	if startDate != "" {
		startTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %v", err)
		}
		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	}

	if endDate != "" {
		endTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, endTime.Location())
	}

	applicationsList, err := s.chargingRepository.GetAllChargingApplications(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get charging applications: %v", err)
	}

	var filteredApplications []domain.ChargingApplication
	for _, chargingApplication := range *applicationsList {
		if (status != "" && chargingApplication.Status != status) || chargingApplication.Status == "deleted" || chargingApplication.Status == "canceled" {
			continue
		}

		if !startTime.IsZero() && chargingApplication.CreatedAt.Before(startTime) {
			continue
		}

		if !endTime.IsZero() && chargingApplication.CreatedAt.After(endTime) {
			continue
		}

		filteredApplications = append(filteredApplications, chargingApplication)
	}

	log.Printf("found %d charging applications with filters", len(filteredApplications))
	return &filteredApplications, nil
}

func (s *ChargingService) GetChargingApplication(ctx context.Context, id uint) (*domain.ChargingApplication, *[]domain.ChargingOrder, error) {
	application, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get charging application: %v", err)
	}
	chargingApplicationOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get charging application orders: %v", err)
	}

	log.Printf("get charging application: %d", id)
	return application, chargingApplicationOrders, nil
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
		UpdatedAt:      time.Now(),
	}

	err = s.chargingRepository.CreateChargingApplication(ctx, newChargingDraft)
	if err != nil {
		return nil, fmt.Errorf("failed to create draft: %v", err)
	}

	log.Printf("created new draft for: %d", creatorId)
	return newChargingDraft, nil
}

func (s *ChargingService) UpdateChargingApplicationPhone(ctx context.Context, phone string, creatorId uint) error {
	application, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return fmt.Errorf("failed to get charging application: %v", err)
	}

	if err := s.validatePhone(phone); err != nil {
		return err
	}

	chargingUpdates := map[string]interface{}{
		"creator_phone": phone,
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, application.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to update phone: %v", err)
	}

	log.Printf("updated draft for: %d", creatorId)
	return nil
}

func (s *ChargingService) FormChargingApplication(ctx context.Context, creatorId uint) error {
	application, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return fmt.Errorf("failed to get charging application: %v", err)
	}

	if application.CreatorPhone == "" {
		return fmt.Errorf("phone number is required to form charging application")
	}

	chargingUpdates := map[string]interface{}{
		"status": "formed",
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, application.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to form application: %v", err)
	}

	log.Printf("formed application for: %d", creatorId)
	return nil
}

func (s *ChargingService) CompleteChargingApplication(ctx context.Context, id uint, moderatorId uint) error {
	application, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %v", err)
	}

	if application.Status != "formed" {
		return fmt.Errorf("can only complete formed applications")
	}

	totalPrice := s.calculateTotalChargingPrice(application) // хелперы?
	// рассчитать для каждого ордера

	chargingUpdates := map[string]interface{}{
		"status":       "completed",
		"moderator_id": moderatorId,
		"completed_at": time.Now(),
		"total_price":  totalPrice,
	}
	err = s.chargingRepository.UpdateChargingApplication(ctx, id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to complete application: %v", err)
	}

	log.Printf("completed application %d by %d", id, moderatorId)
	return nil
}

func (s *ChargingService) RejectChargingApplication(ctx context.Context, id uint, moderatorId uint) error {
	application, err := s.chargingRepository.GetChargingApplicationById(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %v", err)
	}

	if application.Status != "formed" {
		return fmt.Errorf("can only complete formed applications")
	}

	chargingUpdates := map[string]interface{}{
		"status":       "rejected",
		"moderator_id": moderatorId,
	}

	err = s.chargingRepository.UpdateChargingApplication(ctx, id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to reject application: %v", err)
	}

	log.Printf("rejected application %d by %d", id, moderatorId)
	return nil
}

func (s *ChargingService) DeleteChargingApplication(ctx context.Context, creatorId uint) error {
	application, err := s.chargingRepository.GetDraftChargingApplicationByCreator(ctx, creatorId)
	if err != nil {
		return fmt.Errorf("failed to get charging application: %v", err)
	}

	err = s.chargingRepository.DeleteChargingApplicationById(ctx, application.Id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %v", err)
	}

	log.Printf("application deleted: %d", application.Id)
	return nil
}

// Вспомогательные методы для расчетов, хелперы?

func (s *ChargingService) calculateTotalChargingPrice(application *domain.ChargingApplication) float32 {
	return 1000
}

func (s *ChargingService) calculateChargingPriceForOrder(application *domain.ChargingApplication) float32 {
	return 1000
}

func (s *ChargingService) validatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number is required")
	}
	if len(phone) < 10 {
		return fmt.Errorf("phone number is too short")
	}
	return nil
}
