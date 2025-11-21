package customejwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
)

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	claims := &JWTCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TOKEN_SECRET))
}
