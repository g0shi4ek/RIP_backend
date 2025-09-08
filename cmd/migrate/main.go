package main

import (
	"github.com/g0shi4ek/RIP_backend/internal/app/dsn"
	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&domain.ChargingTariff{},
		&domain.User{},
		&domain.ChargingApplication{},
		&domain.ChargingOrder{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}