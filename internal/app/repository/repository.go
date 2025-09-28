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
	rc domain.ISessionStorage
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

	redisClient, err := database.NewRedisClient()
	if err != nil{
		return nil, fmt.Errorf("error initializing redis client: %v", err)
	}

	return &ChargingRepository{
		db: postgresClient,
		mc: minioClient,
		rc: redisClient,
	}, nil
}
