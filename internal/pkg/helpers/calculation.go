package helpers

import (
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

const (
	efficiency          = 0.9 // КПД 90%
	nightTimeMultiplier = 1.3 // Коэффициент для ночного времени
	nightStartHour      = 22  // Начало ночного времени (22:00)
	nightEndHour        = 6   // Конец ночного времени (06:00)
)

func CalculateChargingPriceForOrder(order *domain.ChargingOrder, tariff *domain.ChargingTariff) (float32, float32, error) {
	// Рассчитываем ёмкость нужного заряда в кВт*ч
	chargeNeededKWh := (100 - float32(order.CurrentPercent)) * order.BatteryCapacity / 100

	// Рассчитываем время зарядки в часах
	chargingTimeHours := chargeNeededKWh / (tariff.Power * efficiency)

	// Определяем коэффициент времени (ночь/день)
	timeMultiplier := getTimeMultiplier(order.StartTime)

	// Рассчитываем стоимость
	totalCost := chargingTimeHours * tariff.PricePerHour * timeMultiplier

	return totalCost, chargingTimeHours, nil
}

func getTimeMultiplier(startTime time.Time) float32 {
	hour := startTime.Hour()
	// Ночное время: с 22:00 до 06:00

	if hour >= nightStartHour || hour < nightEndHour {
		return nightTimeMultiplier
	}
	return 1.0 // Дневное время
}
