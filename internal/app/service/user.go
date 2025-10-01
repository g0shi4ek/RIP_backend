package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (s *ChargingService) RegisterChargingUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := helpers.ValidateUser(user); err != nil {
		return nil, err
	}

	existingUser, err := s.chargingRepository.GetChargingUserByLogin(ctx, user.Login)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with login %s already exists", user.Login)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)

	err = s.chargingRepository.CreateChargingUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	user.Password = ""
	log.Printf("user registered: %d, %s", user.Id, user.Login)
	return user, nil
}

func (s *ChargingService) GetChargingUserProfile(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.chargingRepository.GetChargingUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	log.Printf("get user: %d", id)
	user.Password = ""
	return user, nil
}

func (s *ChargingService) UpdateChargingUserProfile(ctx context.Context, id uint, user *domain.User) (*domain.User, error) {
	if err := helpers.ValidateUserLogin(user); err != nil {
		return nil, err
	}

	chargingUser, err := s.chargingRepository.GetChargingUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	if user.Login != "" && user.Login != chargingUser.Login {
		userWithSameLogin, err := s.chargingRepository.GetChargingUserByLogin(ctx, user.Login)
		if err == nil && userWithSameLogin != nil && userWithSameLogin.Id != id {
			return nil, fmt.Errorf("login %s is already used", user.Login)
		}
		chargingUser.Login = user.Login
	}

	chargingUpdates := map[string]interface{}{
		"login":        chargingUser.Login,
	}
	if user.Password != "" {
		if err := helpers.ValidateUserPassword(user); err != nil {
			return nil, err
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		chargingUser.Password = string(hashedPassword)
		chargingUpdates = map[string]interface{}{
			"login":        chargingUser.Login,
			"password":     chargingUser.Password,
		}
	}

	err = s.chargingRepository.UpdateChargingUser(ctx, chargingUser.Id, chargingUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	newUser, err := s.chargingRepository.GetChargingUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	log.Printf("user updated: %d", id)
	return newUser, nil
}

func (s *ChargingService) LoginChargingUser(ctx context.Context, login, password string) (string, error) {
	existingUser, err := s.chargingRepository.GetChargingUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.CreateNewJwtToken(existingUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	log.Printf("user logged in successfully: %d, login: %s", existingUser.Id, existingUser.Login)

	return token, nil
}

func (s *ChargingService) LogoutChargingUser(ctx context.Context, token string) error {
	log.Println("user token", token)
	err := s.chargingRepository.AddTokenToBlacklist(ctx, token, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %v", err)
	}

	log.Printf("user logged out")
	return nil
}

func (s * ChargingService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error){
	return s.chargingRepository.IsTokenBlacklisted(ctx, token)
}