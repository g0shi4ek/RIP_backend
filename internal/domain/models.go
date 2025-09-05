package domain

import (
	"time"
)
type ChargingOrder struct{
	Id int
	TariffId int 
	BatteryCapacity float32
	CurrentPercent int
	StartTime time.Time // if night => price > 
	OrderPrice float32 
	CreatedAt time.Time
	Deleted bool
}

type ChargingTariff struct {
	Id int
	NameofTariff string
	Description string
	ImageUrl string
	PricePerHour float32
	Power float32
	CreatedAt time.Time
	Deleted bool
}

type ChargingApplication struct{
	Id int
	Price float32
	AmountOfOrders int
	OrdersList *[]ChargingOrder
	Status int // 1 - что-то добавлено, 0 - обработана
	CreatedAt time.Time
	Deleted bool
}