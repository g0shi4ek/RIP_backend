package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"gorm.io/gorm"
)

func (r *ChargingRepository) CreateChargingUser(ctx context.Context, user *domain.User) error {    
    err := r.db.WithContext(ctx).Create(user).Error
    if err != nil {
        return fmt.Errorf("failed to create user: %v", err)
    }
    
    log.Printf("user created: %d, login: %s", user.Id, user.Login)
    return nil
}

func (r *ChargingRepository) UpdateChargingUser(ctx context.Context, user *domain.User) error { // накинуть в бд уникальность на пароль    
    err := r.db.WithContext(ctx).Model(&domain.User{}).
        Where("id = ?", user.Id).
        Updates(map[string]interface{}{
            "login":        user.Login,
            "password":     user.Password,
            "is_moderator": user.IsModerator,
        }).Error

    if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("user with id %d not found", user.Id)
        }
        return fmt.Errorf("failed to update user: %v", err)
    }

    log.Printf("user updated: %d", user.Id)
    return nil
}

func (r *ChargingRepository) GetChargingUserById(ctx context.Context, id uint) (*domain.User, error) {
    var user domain.User
    
    err := r.db.WithContext(ctx).
        Model(&domain.User{}).
        Where("id = ?", id).
        First(&user).Error
        
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with id %d not found", id)
        }
        return nil, fmt.Errorf("failed to get user: %v", err)
    }
    
    return &user, nil
}

func (r *ChargingRepository) GetChargingUserByLogin(ctx context.Context, login string) (*domain.User, error) {
    var user domain.User
    
    err := r.db.WithContext(ctx).
        Model(&domain.User{}).
        Where("login = ?", login).
        First(&user).Error
        
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user with login %s not found", login)
        }
        return nil, fmt.Errorf("failed to get user by login: %v", err)
    }
    
    return &user, nil
}