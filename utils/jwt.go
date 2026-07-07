package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/omnlgy/jadwalin/internal/config"
	"github.com/omnlgy/jadwalin/internal/domain"
)

type JWTClaims struct {
	UserID      string
	Role        string
	PhoneNumber string
	jwt.RegisteredClaims
}

func GenerateJWT(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID:      user.ID.String(),
		Role:        string(user.Role),
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT_SECRET))
}

func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidJWT
	}

	return claims, nil
}
