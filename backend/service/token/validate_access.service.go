package token

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

// ValidateAccessToken validates and parses an access token
func (s *tokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
