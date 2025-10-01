package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

func (h *ChargingHandler) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	createdUser, err := h.chargingService.RegisterChargingUser(ctx, &user)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdUser.ToResponse())
}

func (h *ChargingHandler) GetUserById(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	user, err := h.chargingService.GetChargingUserProfile(ctx, id)
	if err != nil {
		h.ErrorHandler(c, fmt.Errorf("failed to get user: %v", err))
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

func (h *ChargingHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	id, err := helpers.ValidateID(c.Param("id"))
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	var userUpdateRequest domain.UserUpdateRequest
	if err := c.ShouldBindJSON(&userUpdateRequest); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	user := domain.User{
		Id:       id,
		Login:    userUpdateRequest.Login,
		Password: userUpdateRequest.Password,
	}

	newUser, err := h.chargingService.UpdateChargingUserProfile(ctx, id, &user)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
		"user":    newUser.ToResponse(),
	})
}

func (h *ChargingHandler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var userRequest domain.UserLoginRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	newUser, err := h.chargingService.LoginChargingUser(ctx, userRequest.Login, userRequest.Password)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user login successful",
		"user":    newUser.ToResponse(),
	})
}

func (h *ChargingHandler) LogOutUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	token := "bearer: ghjkl"

	err := h.chargingService.LogoutChargingUser(ctx, token)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
