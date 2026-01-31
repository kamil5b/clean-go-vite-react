package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// GenerateRefreshToken generates a new refresh token
func (s *tokenService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.RefreshTokenExpiry).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "go-vite-react",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.RefreshTokenSecret))
}
