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
	for _, tariff := range *tariffList {
		if tariff.Id == id {
			targetTariff = tariff
		}
	}

	return &targetTariff, nil
}

func (r *ChargingRepository) GetChargingApplicationById(ctx context.Context, id int) (*domain.ChargingApplication, error) {
	application := domain.ChargingApplication{
		Id:             1,
		Price:          5400.3,
		AmountOfOrders: 2,
		Status:         "draft",
		CreatedAt:      time.Now(),
		IsDeleted:      false,
	}

	return &application, nil
}

func (r *ChargingRepository) GetAllTariffs(ctx context.Context) (*[]domain.ChargingTariff, error) {
	tariffsList := []domain.ChargingTariff{
		{
			Id:           1,
			NameofTariff: "Быстрая зарядка DC (будни)",
			Description:  "Зарядка постоянным током 50-150 кВт",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 12.0,  // цена за кВт*ч
			Power:        150.0, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
		{
			Id:           2,
			NameofTariff: "Быстрая зарядка DC (выходные)",
			Description:  "Зарядка постоянным током 50-150 кВт",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 15.0,  // цена за кВт*ч
			Power:        150.0, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
		{
			Id:           3,
			NameofTariff: "AC зарядка Level 2 (будни)",
			Description:  "Зарядка переменным током 22 кВт",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 8.0,  // цена за кВт*ч
			Power:        22.0, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
		{
			Id:           4,
			NameofTariff: "AC зарядка Level 2 (выходные)",
			Description:  "Зарядка переменным током 22 кВт",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 10.0, // цена за кВт*ч
			Power:        22.0, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
		{
			Id:           5,
			NameofTariff: "Медленная зарядка Level 1",
			Description:  "Домашняя зарядка 3.7-7.4 кВт",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 5.0, // цена за кВт*ч
			Power:        7.4, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
		{
			Id:           6,
			NameofTariff: "Ультрабыстрая зарядка DC",
			Description:  "Зарядка 350 кВт (Tesla Supercharger)",
			ImageUrl:     "http://127.0.0.1:9000/charging-images/image.png",
			PricePerHour: 18.0,  // цена за кВт*ч
			Power:        350.0, // кВт
			CreatedAt:    time.Now(),
			IsDeleted:    false,
		},
	}

	if len(tariffsList) == 0 {
		return nil, fmt.Errorf("there are no tariffs")
	}
	log.Println("tariffList founded")
	return &tariffsList, nil
}

func (r *ChargingRepository) SearchTariffs(ctx context.Context, tariffQuery string) (*[]domain.ChargingTariff, error) {
	tariffQuery = strings.ToLower(tariffQuery)
	var results []domain.ChargingTariff
	tariffs, _ := r.GetAllTariffs(ctx)

	for _, tariff := range *tariffs {
		if strings.Contains(strings.ToLower(tariff.NameofTariff), tariffQuery) ||
			strings.Contains(strings.ToLower(tariff.Description), tariffQuery) {
			results = append(results, tariff)
		}
	}

	return &results, nil
}
