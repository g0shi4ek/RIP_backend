package jwt

import (
	"log"
	"os"
	"time"

	"github.com/g0shi4ek/RIP_backend/internal/domain"
	"github.com/golang-jwt/jwt"
)

func setRole(user *domain.User) string {
	if user.IsModerator {
		return "moderator"
	}
	return "client"
}

func CreateNewJwtToken(user *domain.User) (string, error) {
	key := os.Getenv("JWT_KEY")
	keyBytes := []byte(key)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,                          // Subject (user identifier)
		"iss": "app",                            // Issuer
		"aud": setRole(user),                    // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(keyBytes)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJwtToken(tokenString string) (jwt.MapClaims, error) {
	key := os.Getenv("JWT_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func ExtractUserDataFromToken(tokenString string) (uint, string, error) {
	claims, err := VerifyJwtToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	//subject
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", jwt.ErrInvalidKey
	}

	//роль (audience)
	aud, ok := claims["aud"].(string)
	if !ok {
		return 0, "", jwt.ErrInvalidKey
	}
	log.Println(sub, aud)

	exp, ok := claims["exp"].(float64)
	if !ok {
		return 0, "", jwt.ErrInvalidKey
	}
	if time.Now().Unix() > int64(exp) {
		return 0, "", jwt.ErrInvalidKey
	}

	return uint(sub), aud, nil
}
