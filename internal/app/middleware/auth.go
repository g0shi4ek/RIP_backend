package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/g0shi4ek/RIP_backend/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	chargingService domain.IChargingService
}

func NewAuthMiddleware(chargingService domain.IChargingService) *AuthMiddleware {
	return &AuthMiddleware{
		chargingService: chargingService,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"description": "authorization header required",
			})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		isBlacklisted, err := m.chargingService.IsTokenBlacklisted(ctx, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":      "error",
				"description": "failed to check token",
			})
			c.Abort()
			return
		}

		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"description": "token is blacklisted",
			})
			c.Abort()
			return
		}

		userId, role, err := jwt.ExtractUserDataFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":      "error",
				"description": "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Set("userRole", role)
		c.Next()
	}
}

func (m *AuthMiddleware) Role(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "user role not found",
			})
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "invalid user role",
			})
			c.Abort()
			return
		}

		hasAccess := false
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{
				"status":      "error",
				"description": "access denied",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetCurrentUserID(c *gin.Context) (uint, error) {
	userId, ok:= c.Get("userId")
	if !ok{
		return 0, fmt.Errorf("user id not found in context")
	}
	fmt.Println(userId)

	id, ok := userId.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	return id, nil
}

func GetCurrentUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get("userRole")
	if !ok {
		return "", fmt.Errorf("user role not found in context")
	}

	roleStr, ok := role.(string)
	if !ok {
		return "", fmt.Errorf("invalid user role type")
	}

	return roleStr, nil
}