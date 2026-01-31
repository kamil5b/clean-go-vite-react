package csrf

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken generates a random CSRF token
func (s *csrfService) GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
