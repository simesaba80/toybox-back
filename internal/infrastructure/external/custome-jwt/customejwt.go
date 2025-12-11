package customejwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
)

func GenerateToken(userID uuid.UUID) (string, error) {
	claims := &schema.JWTCustomClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.TOKEN_SECRET))
}
