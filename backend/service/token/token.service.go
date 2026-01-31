package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TokenConfig holds JWT configuration
type TokenConfig struct {
	AccessTokenSecret  string
	AccessTokenExpiry  time.Duration
	RefreshTokenSecret string
	RefreshTokenExpiry time.Duration
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	jwt.StandardClaims
}

// TokenService handles token generation and validation
type TokenService interface {
	GenerateAccessToken(userID uuid.UUID, email, name string) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateAccessToken(tokenString string) (*TokenClaims, error)
	ValidateRefreshToken(tokenString string) (*TokenClaims, error)
}

// tokenService implements TokenService
type tokenService struct {
	config TokenConfig
}

// NewTokenService creates a new token service
func NewTokenService(config TokenConfig) TokenService {
	return &tokenService{
		config: config,
	}
}
