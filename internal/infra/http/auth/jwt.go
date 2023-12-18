package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jwtConfig struct {
	secretKey          []byte
	expirationDuration time.Duration
}

var config jwtConfig

func SetJWTConfig(secretKey []byte, expDur time.Duration) {
	config = jwtConfig{secretKey: secretKey, expirationDuration: expDur}
}

func CreateToken(userID uint64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"iss":      "basket-service/jwt.go",
		"exp":      time.Now().Add(config.expirationDuration).Unix(),
		"aud":      "user_" + username,
	})

	tokenString, err := token.SignedString(config.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return config.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	return claims, nil
}
