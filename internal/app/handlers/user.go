package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register user
// @Description Register new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.UserRequest true "User credentials"
// @Success 201 {integer} integer "Created id"
// @Failure 400 {object} object "Bad request"
// @Failure 409 {object} object "Conflict - already exists"
// @Failure 500 {object} object "Internal server error"
// @Router /users/register [post]
func (h *ChargingHandler) RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var request domain.UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	user := domain.User{
		Login:    request.Login,
		Password: request.Password,
	}

	createdUser, err := h.chargingService.RegisterChargingUser(ctx, &user)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdUser.ToResponse())
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Get user profile by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 404 {object} object "User not found"
// @Failure 500 {object} object "Internal server error"
// @Router /users/{id} [get]
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

// UpdateUser godoc
// @Summary Update user
// @Description Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param request body domain.UserRequest true "User data"
// @Security BearerAuth
// @Success 200 {object} object "User updated successfully"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 404 {object} object "User not found"
// @Failure 409 {object} object "Conflict - already exists"
// @Failure 500 {object} object "Internal server error"
// @Router /users/{id} [put]
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

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user and get JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.UserRequest true "User credentials"
// @Success 200 {object} object "Login response with token"
// @Failure 400 {object} object "Bad request"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /users/login [post]
func (h *ChargingHandler) LoginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var userRequest domain.UserLoginRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.ErrorHandler(c, fmt.Errorf("%w: %v", ErrInvalidRequestBody, err))
		return
	}

	token, err := h.chargingService.LoginChargingUser(ctx, userRequest.Login, userRequest.Password)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

// LogOutUser godoc
// @Summary Logout user
// @Description Invalidate user's JWT token
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Logout successful"
// @Failure 401 {object} object "Unauthorized"
// @Failure 500 {object} object "Internal server error"
// @Router /users/logout [post]
func (h *ChargingHandler) LogOutUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	err := h.chargingService.LogoutChargingUser(ctx, token)
	if err != nil {
		h.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
