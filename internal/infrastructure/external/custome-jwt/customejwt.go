package customejwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
)

func GenerateToken(userID string) (string, error) {
	claims := &schema.JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TOKEN_SECRET))
}

func RegenerateToken(refreshToken string) (string, error) {
	claims := &schema.JWTCustomClaims{
		RefreshToken: refreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TOKEN_SECRET))
}
