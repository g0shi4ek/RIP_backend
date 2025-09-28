package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

func (r *ChargingRepository) CreateChargingUser(ctx context.Context, user *domain.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	
	log.Printf("repo: user created, %d", user.Id)
	return nil
}

func (r *ChargingRepository) UpdateChargingUser(ctx context.Context, id uint, chargingUpdates map[string]interface{}) error {
	if len(chargingUpdates) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(chargingUpdates).Error

	if err != nil {
		return fmt.Errorf("failed to update charging user: %v", err)
	}
	log.Printf("repo: user updated, %d", id)
	return nil
}

func (r *ChargingRepository) GetChargingUserById(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	log.Printf("repo: user retrieved, %d", id)
	return &user, nil
}

func (r *ChargingRepository) GetChargingUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("login = ?", login).
		First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user by login: %v", err)
	}
	log.Printf("repo: user retrieved, %s", login)
	return &user, nil
}

func (r *ChargingRepository) AddTokenToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	return r.rc.AddToBlacklist(ctx, token, ttl)
}

func (r *ChargingRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return r.rc.IsInBlacklist(ctx, token)
}