package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

type ChargingRepository struct{}

func NewChargingRepository() (*ChargingRepository, error) {
	return &ChargingRepository{}, nil
}

func (r *ChargingRepository) GetTariffById(ctx context.Context, id int) (*domain.ChargingTariff, error) {

	var targetTariff domain.ChargingTariff
	tariffList, _ := r.GetAllTariffs(ctx)
	for _, tariff  := range *tariffList{
		if tariff.Id == id{
			targetTariff = tariff
		}
	}

	return &targetTariff, nil
}

func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, id int) (*domain.ChargingApplication, error){
	orders := []domain.ChargingOrder{
        {
            Id:              1,
            TariffId:        1,
            BatteryCapacity: 75.5,
            CurrentPercent:  45,
            StartTime:       time.Now(),
            OrderPrice:      1500.0,
            CreatedAt:       time.Now(),
            Deleted:         false,
        },
        {
            Id:              2,
            TariffId:        2,
            BatteryCapacity: 60.0,
            CurrentPercent:  30,
            StartTime:       time.Now(),
            OrderPrice:      1200.0,
            CreatedAt:       time.Now().Add(-2 * time.Hour),
            Deleted:         false,
        },
    }

	application := domain.ChargingApplication{
        Id:            1,
        Price:         5400.3,
        AmountOfOrders: 2,
        OrdersList:    &orders,
        Status:        1,
        CreatedAt:     time.Now(),
        Deleted:       false,
    }

	return &application, nil
}

func (r *ChargingRepository) GetAllTariffs(ctx context.Context) (*[]domain.ChargingTariff, error){
	tariffsList := []domain.ChargingTariff{
		{
			Id:           1,
			NameofTariff: "Быстрый тариф в будни",
			Description:  "Зарядка уровня 1",
			ImageUrl:     "/resources/img/image.png",
			PricePerHour: 5.0,
			Power:        17.6,
			CreatedAt:    time.Now(),
			Deleted:      false,
		},
		{
			Id:           2,
			NameofTariff: "Быстрый тариф в выходные",
			Description:  "Зарядка уровня 1",
			ImageUrl:     "/resources/img/image.png",
			PricePerHour: 6.0,
			Power:        17.6,
			CreatedAt:    time.Now(),
			Deleted:      false,
		},
		{
			Id:           3,
			NameofTariff: "Медленный тариф в будни",
			Description:  "Зарядка уровня 2",
			ImageUrl:     "/resources/img/image.png",
			PricePerHour: 6.0,
			Power:        10.1,
			CreatedAt:    time.Now(),
			Deleted:      false,
		},
		{
			Id:           4,
			NameofTariff: "Медленный тариф в выходные",
			Description:  "Зарядка уровня 2",
			ImageUrl:     "/resources/img/image.png",
			PricePerHour: 6.0,
			Power:        10.1,
			CreatedAt:    time.Now(),
			Deleted:      false,
		},
    }
	
	if len(tariffsList) == 0{
		return nil, fmt.Errorf("there are no tariffs")
	}
	log.Println("tariffList founded")
	return &tariffsList, nil
}


func (r *ChargingRepository) SearchTariffs(ctx context.Context, query string) (*[]domain.ChargingTariff, error) {
    query = strings.ToLower(query)
    var results []domain.ChargingTariff
	tariffs, _ := r.GetAllTariffs(ctx)
	
    for _, tariff := range *tariffs {
        if strings.Contains(strings.ToLower(tariff.NameofTariff), query) ||
           strings.Contains(strings.ToLower(tariff.Description), query) {
            results = append(results, tariff)
        }
    }
    
    return &results, nil
}