package domain

import (
	"time"
)

// M-M
type ChargingOrder struct {
	Id              uint                `gorm:"primaryKey" json:"id"`
	ApplicationId   uint                `gorm:"not null;uniqueIndex:idx_order_unique" json:"application_id"`
	TariffId        uint                `gorm:"not null;uniqueIndex:idx_order_unique" json:"tariff_id"`
	Application     ChargingApplication `gorm:"foreignKey:ApplicationId" json:"-"`
	Tariff          ChargingTariff      `gorm:"foreignKey:TariffId" json:"tariff"`
	BatteryCapacity float32             `gorm:"not null" json:"battery_capacity"`
	CurrentPercent  int                 `gorm:"not null" json:"current_percent"`
	StartTime       time.Time           `gorm:"not null" json:"start_time"` // if night => price >
	EstimatedTime   float32             `json:"estimated_time"`             // расчетное время зарядки в часах
	CalculatedPrice float32             `json:"calculated_price"`           // расчетная стоимость для этой услуги
}

// Service
type ChargingTariff struct {
	Id           uint    `gorm:"primaryKey" json:"id"`
	NameofTariff string  `gorm:"type:varchar(50);not null" json:"nameof_tariff"`
	Description  string  `gorm:"type:varchar(100);not null" json:"description"`
	ImageUrl     string  `gorm:"type:varchar(100);null" json:"image_url,omitempty"`
	PricePerHour float32 `gorm:"not null" json:"price_per_hour"`
	Power        float32 `gorm:"not null" json:"power"` // мощность зарядки в кВт
	IsDeleted    bool    `gorm:"default:false" json:"is_deleted"`
}

// Application
type ChargingApplication struct {
	Id             uint            `gorm:"primaryKey" json:"id"`
	TotalPrice     float32         `json:"total_price"`
	CreatorId      uint            `gorm:"not null" json:"creator_id"`
	ModeratorId    uint            `gorm:"default:null" json:"moderator_id,omitempty"`
	Creator        User            `gorm:"foreignKey:CreatorId" json:"-"`
	Moderator      User            `gorm:"foreignKey:ModeratorId" json:"-"`
	CreatorPhone   string          `json:"creator_phone,omitempty"`
	AmountOfOrders uint            `gorm:"default:0" json:"amount_of_orders"`
	Status         string          `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, deleted, formed, completed, rejected
	CreatedAt      time.Time       `gorm:"autoCreateTime" json:"created_at"`
	FormedAt       time.Time       `gorm:"default:null" json:"formed_at,omitempty"`
	CompletedAt    time.Time       `gorm:"default:null" json:"completed_at,omitempty"`
	Orders         []ChargingOrder `gorm:"foreignKey:ApplicationId" json:"orders,omitempty"`
}

type User struct {
	Id          uint   `gorm:"primary_key" json:"id"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"password"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}
