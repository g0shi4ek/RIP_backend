package domain

import "time"

type ChargingOrderResponse struct {
	TariffId        uint           `json:"tariff_id"`
	Tariff          TariffResponse `json:"tariff"`
	BatteryCapacity float32        `json:"battery_capacity"`
	CurrentPercent  int            `json:"current_percent"`
	StartTime       time.Time      `json:"start_time"`
	EstimatedTime   float32        `json:"estimated_time"`
	CalculatedPrice float32        `json:"calculated_price"`
}

type TariffResponse struct {
	Id           uint    `json:"id"`
	NameofTariff string  `json:"nameof_tariff"`
	Description  string  `json:"description"`
	ImageUrl     string  `json:"image_url,omitempty"`
	PricePerHour float32 `json:"price_per_hour"`
	Power        float32 `json:"power"`
}

type ChargingApplicationResponse struct {
	Id             uint                    `json:"id"`
	TotalPrice     float32                 `json:"total_price"`
	CreatorId      uint                    `json:"creator_id"`
	ModeratorId    uint                    `json:"moderator_id,omitempty"`
	CreatorPhone   string                  `json:"creator_phone,omitempty"`
	AmountOfOrders uint                    `json:"amount_of_orders"`
}

type ChargingDraftResponse struct {
	Id             uint                    `json:"id"`
	AmountOfOrders uint                    `json:"amount_of_orders"`
}

type UserResponse struct {
	Id          uint   `json:"id"`
	Login       string `json:"login"`
}

func (o *ChargingOrder) ToResponse() ChargingOrderResponse {
	return ChargingOrderResponse{
		TariffId:        o.TariffId,
		Tariff:          o.Tariff.ToResponse(),
		BatteryCapacity: o.BatteryCapacity,
		CurrentPercent:  o.CurrentPercent,
		StartTime:       o.StartTime,
		EstimatedTime:   o.EstimatedTime,
		CalculatedPrice: o.CalculatedPrice,
	}
}

func (t *ChargingTariff) ToResponse() TariffResponse {
	return TariffResponse{
		Id:           t.Id,
		NameofTariff: t.NameofTariff,
		Description:  t.Description,
		ImageUrl:     t.ImageUrl,
		PricePerHour: t.PricePerHour,
		Power:        t.Power,
	}
}

func (a *ChargingApplication) ToResponse() ChargingApplicationResponse {
	return ChargingApplicationResponse{
		Id:             a.Id,
		TotalPrice:     a.TotalPrice,
		CreatorId:      a.CreatorId,
		ModeratorId:    a.ModeratorId,
		CreatorPhone:   a.CreatorPhone,
		AmountOfOrders: a.AmountOfOrders,
	}
}

func (a *ChargingApplication) ToDraftResponse() ChargingDraftResponse {
	return ChargingDraftResponse{
		Id:             a.Id,
		AmountOfOrders: a.AmountOfOrders,
	}
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		Id:          u.Id,
		Login:       u.Login,
	}
}
