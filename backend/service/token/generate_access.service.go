package token

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

// GenerateAccessToken generates a new access token
func (s *tokenService) GenerateAccessToken(userID uuid.UUID, email, name string) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.AccessTokenExpiry).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "go-vite-react",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.AccessTokenSecret))
}
