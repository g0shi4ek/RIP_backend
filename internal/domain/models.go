package domain

import (
	"time"
)

// M-M
type ChargingOrder struct {
	Id              int
	TariffId        int
	ApplicationId   int
	BatteryCapacity float32
	CurrentPercent  int
	StartTime       time.Time // if night => price >
	OrderPrice      float32
	CreatedAt       time.Time
	IsDeleted       bool
}

// Service
type ChargingTariff struct {
	Id           int
	NameofTariff string
	Description  string
	ImageUrl     string
	PricePerHour float32
	Power        float32
	CreatedAt    time.Time
	IsDeleted    bool
}

// Application
type ChargingApplication struct {
	Id             int
	Price          float32
	AmountOfOrders int
	Status         string
	CreatedAt      time.Time
	IsDeleted      bool
}
