package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
		if userRole == "client" && chargingApplication.CreatorId != creatorId {
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

	chargingOrders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, chargingApplication.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for application: %v", err)
	}

	for _, order := range *chargingOrders {
		err := s.callAsyncCalculationService(order)
		if err != nil {
			log.Printf("Failed to start async calculation for order app=%d tariff=%d: %v",
				order.ApplicationId, order.TariffId, err)
		}
	}

	chargingUpdates := map[string]interface{}{
		"status":       "completed",
		"moderator_id": moderatorId,
		"completed_at": time.Now(),
		"total_price":  0, //  пока расчеты не завершены
	}

	err = s.chargingRepository.UpdateChargingApplication(ctx, id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to complete application: %v", err)
	}

	respApplication, _ := s.chargingRepository.GetChargingApplicationById(ctx, id)
	log.Printf("completed application %d by %d - async calculations started", id, moderatorId)
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

func (s *ChargingService) UpdateApplicationTotalPrice(ctx context.Context, applicationId uint) error {
	orders, err := s.chargingRepository.GetChargingOrdersByApplicationId(ctx, applicationId)
	if err != nil {
		return fmt.Errorf("failed to get orders: %v", err)
	}

	var totalPrice float32
	var completedCalculations int

	for _, order := range *orders {
		totalPrice += order.CalculatedPrice
		completedCalculations++
	}

	updates := map[string]interface{}{
		"total_price": totalPrice,
	}

	err = s.chargingRepository.UpdateChargingApplication(ctx, applicationId, updates)
	if err != nil {
		return fmt.Errorf("failed to update total price: %v", err)
	}

	log.Printf("Updated total price for application %d: %.2f (%d/%d calculations completed)",
		applicationId, totalPrice, completedCalculations, len(*orders))

	return nil
}

func (s *ChargingService) callAsyncCalculationService(order domain.ChargingOrder) error {
	requestData := map[string]interface{}{
		"application_id":   order.ApplicationId,
		"tariff_id":        order.TariffId,
		"battery_capacity": order.BatteryCapacity,
		"current_percent":  order.CurrentPercent,
		"tariff_power":     order.Tariff.Power,
		"price_per_hour":   order.Tariff.PricePerHour,
		"start_time":       order.StartTime.Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %v", err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(
		"http://localhost:8000/api/charging_calculate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to call Django service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("django service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Failed to decode Django response: %v", err)
	}

	log.Printf("Successfully started async calculation for order app=%d tariff=%d. Response: %s",
		order.ApplicationId, order.TariffId, response.Message)

	return nil
}
