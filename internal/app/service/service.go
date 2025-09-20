package service

import (
	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

type ChargingService struct {
	chargingRepository domain.IChargingRepository
}

func NewChargingService(repo domain.IChargingRepository) (*ChargingService, error) {
	return &ChargingService{
		chargingRepository: repo,
	}, nil
}