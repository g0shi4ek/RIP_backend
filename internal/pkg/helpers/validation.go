package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
)

var (
	ErrInvalidID          = fmt.Errorf("invalid ID")
	ErrInvalidRequestBody = fmt.Errorf("invalid request body")
	ErrInvalidPhone       = fmt.Errorf("invalid phone number")
	ErrImageRequired      = fmt.Errorf("image file is required")
	ErrInvalidTariffData  = fmt.Errorf("invalid tariff data")
	ErrInvalidLogin       = fmt.Errorf("invalid user login")
	ErrInvalidPassword    = fmt.Errorf("invalid user password")
	ErrInvalidOrderData   = fmt.Errorf("invalid battery capaciry or current percent")
	ErrImageTooLarge      = fmt.Errorf("image size too large, maximum 5MB")
	ErrInvalidImageType   = fmt.Errorf("only image files are allowed")
	ErrFailedReadImage    = fmt.Errorf("failed to read image file")
	ErrInvalidDate        = fmt.Errorf("invalid date format")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
	ErrForbidden          = fmt.Errorf("forbidden")
)

func ValidatePhone(phone string) error {
	if phone == "" {
		return ErrInvalidPhone
	}
	if len(phone) < 10 {
		return ErrInvalidPhone
	}
	return nil
}

func ValidateChargingOrder(chargingOrder *domain.ChargingOrder) error {
	if chargingOrder.BatteryCapacity <= 0 {
		return ErrInvalidOrderData
	}
	if chargingOrder.CurrentPercent < 0 {
		return ErrInvalidOrderData
	}
	return nil
}

func ValidateTariff(tariff *domain.ChargingTariff) error {
	if tariff.NameofTariff == "" {
		return ErrInvalidTariffData
	}
	if tariff.Description == "" {
		return ErrInvalidTariffData
	}
	if tariff.PricePerHour <= 0 {
		return ErrInvalidTariffData
	}
	if tariff.Power <= 0 {
		return ErrInvalidTariffData
	}
	return nil
}

func ValidateUser(user *domain.User) error {
	if strings.TrimSpace(user.Login) == "" {
		return ErrInvalidLogin
	}
	if len(user.Login) < 3 {
		return ErrInvalidLogin
	}
	if strings.TrimSpace(user.Password) == "" {
		return ErrInvalidPassword
	}
	if len(user.Password) < 3 {
		return ErrInvalidPassword
	}
	return nil
}

func ValidateID(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, ErrInvalidID
	}
	return uint(id), nil
}

func ValidateImage(fileHeader *multipart.FileHeader) ([]byte, error) {
	if fileHeader.Size > 5<<20 {
		return nil, ErrImageTooLarge
	}

	if !strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
		return nil, ErrInvalidImageType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, ErrFailedReadImage
	}
	defer file.Close()

	tariffImageData, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrFailedReadImage
	}

	return tariffImageData, nil
}
