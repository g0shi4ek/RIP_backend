package repository

import (
	"fmt"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/database"
	"gorm.io/gorm"
)

type ChargingRepository struct {
	db *gorm.DB
	mc domain.ITariffImagesStorage
}

func NewChargingRepository() (*ChargingRepository, error) {
	postgresClient, err := database.NewPostgresClient()
	if err != nil {
		return nil, fmt.Errorf("error initializing postgres client: %v", err)
	}

	minioClient, err := database.NewMinioClient()
	if err != nil {
		return nil, fmt.Errorf("error initializing minio client: %v", err)
	}

	return &ChargingRepository{
		db: postgresClient,
		mc: minioClient,
	}, nil
}
