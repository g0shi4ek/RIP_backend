package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *ChargingService) RegisterChargingUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := s.validateUser(user); err != nil {
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

	log.Printf("user registered: %d, %s", user.Id, user.Login)
	return user, nil
}

func (s *ChargingService) GetChargingUserProfile(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.chargingRepository.GetChargingUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	log.Printf("get user: %d", id)
	return user, nil
}

func (s *ChargingService) UpdateChargingUserProfile(ctx context.Context, id uint, user *domain.User) error {
	if err := s.validateUser(user); err != nil {
		return err
	}

	existingUser, err := s.chargingRepository.GetChargingUserById(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if user.Login != "" && user.Login != existingUser.Login {
		userWithSameLogin, err := s.chargingRepository.GetChargingUserByLogin(ctx, user.Login)
		if err == nil && userWithSameLogin != nil && userWithSameLogin.Id != id {
			return fmt.Errorf("login %s is already used", user.Login)
		}
		existingUser.Login = user.Login
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		existingUser.Password = string(hashedPassword)
	}

	chargingUpdates := map[string]interface{}{
		"login":        existingUser.Login,
		"password":     existingUser.Password,
		"is_moderator": user.IsModerator,
	}

	err = s.chargingRepository.UpdateChargingUser(ctx, existingUser.Id, chargingUpdates)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	log.Printf("user updated: %d", id)
	return nil
}

func (s *ChargingService) LoginChargingUser(ctx context.Context, login, password string) (*domain.User, error) {
	existingUser, err := s.chargingRepository.GetChargingUserByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	log.Printf("user logged in successfully: %d, login: %s", existingUser.Id, existingUser.Login)
	return existingUser, nil
}

func (s *ChargingService) LogoutChargingUser(ctx context.Context, token string) error {
	log.Printf("user logged out")
	return nil
}

func (s *ChargingService) validateUser(user *domain.User) error { // в хелперы?
	if strings.TrimSpace(user.Login) == "" {
		return fmt.Errorf("login is required")
	}
	if len(user.Login) < 3 {
		return fmt.Errorf("login must be at least 3 characters long")
	}
	if strings.TrimSpace(user.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(user.Password) < 3 {
		return fmt.Errorf("password must be at least 3 characters long")
	}
	return nil
}
