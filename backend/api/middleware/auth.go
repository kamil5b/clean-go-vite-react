package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TokenValidator validates JWT tokens
type TokenValidator interface {
	ValidateToken(tokenString string) (uuid.UUID, error)
}

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(validator TokenValidator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try to get token from cookie first
			token := ""
			cookie, err := c.Cookie("access_token")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
			}

			// If no cookie, try Authorization header
			if token == "" {
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader != "" {
					// Expected format: "Bearer <token>"
					parts := strings.Split(authHeader, " ")
					if len(parts) == 2 && parts[0] == "Bearer" {
						token = parts[1]
					}
				}
			}

			// No token found
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authentication token",
				})
			}

			// Validate token
			userID, err := validator.ValidateToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			// Store user ID in context
			c.Set("user_id", userID)

			return next(c)
		}
	}
}
