package domain

import (
	"time"
)

// M-M
type ChargingOrder struct {
	Id uint `gorm:"primaryKey"`
	// Поля для связей + составной уникальный индекс
	ApplicationId uint `gorm:"not null;uniqueIndex:idx_order_unique"`
	TariffId      uint `gorm:"not null;uniqueIndex:idx_order_unique"`

	// Внешние ключи
	Application ChargingApplication `gorm:"foreignKey:ApplicationId"`
	Tariff      ChargingTariff      `gorm:"foreignKey:TariffId"`

	BatteryCapacity float32   `gorm:"not null"`
	CurrentPercent  int       `gorm:"not null"`
	StartTime       time.Time `gorm:"not null"` // if night => price >
	EstimatedTime   float32   // расчетное время зарядки в часах
	CalculatedPrice float32   // расчетная стоимость для этой услуги

	IsDeleted bool `gorm:"default:false"`
}

// Service
type ChargingTariff struct {
	Id           uint    `gorm:"primaryKey"`
	NameofTariff string  `gorm:"type:varchar(50);not null"`
	Description  string  `gorm:"type:varchar(100);not null"`
	ImageUrl     string  `gorm:"type:varchar(100);null"`
	PricePerHour float32 `gorm:"not null"`
	Power        float32 `gorm:"not null"` // мощность зарядки в кВт
	IsDeleted    bool    `gorm:"default:false"`
}

// Application
type ChargingApplication struct {
	Id          uint `gorm:"primaryKey"`
	TotalPrice  float32
	CreatorId   uint `gorm:"not null"`
	ModeratorId uint

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`

	CreatorPhone   string
	AmountOfOrders uint      `gorm:"default:0"`
	Status         string    `gorm:"type:varchar(20);not null;default:'draft'"` // черновик, удалён, сформирован, завершён, отклонён
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CompletedAt    time.Time `gorm:"default:null"`
}

type User struct {
	Id          uint   `gorm:"primary_key" json:"id"`
	Login       string `gorm:"type:varchar(25);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}
