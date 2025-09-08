package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ChargingRepository struct {
	db *gorm.DB
}

func NewChargingRepository(dsn string) (*ChargingRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &ChargingRepository{
		db: db,
	}, nil
}
