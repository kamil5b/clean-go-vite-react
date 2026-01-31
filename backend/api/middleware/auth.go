package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
	"github.com/labstack/echo/v4"
)

const (
	AccessTokenCookie = "access_token"
	RefreshTokenCookie = "refresh_token"
	UserIDCtxKey      = "user_id"
	UserEmailCtxKey   = "user_email"
	ClaimsCtxKey      = "claims"
)

// AuthMiddleware validates JWT token from HTTP-only cookie
func AuthMiddleware(tokenService token.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from HTTP-only cookie
			cookie, err := c.Cookie(AccessTokenCookie)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authentication token",
				})
			}

			// Validate token
			claims, err := tokenService.ValidateAccessToken(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			// Store user info in context
			c.Set(UserIDCtxKey, claims.UserID.String())
			c.Set(UserEmailCtxKey, claims.Email)
			c.Set(ClaimsCtxKey, claims)

			return next(c)
		}
	}
}

// OptionalAuthMiddleware validates JWT but doesn't require it
func OptionalAuthMiddleware(tokenService token.TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Try to get token from HTTP-only cookie
			cookie, err := c.Cookie(AccessTokenCookie)
			if err != nil {
				// Token not present, continue without authentication
				return next(c)
			}

			// Validate token
			claims, err := tokenService.ValidateAccessToken(cookie.Value)
			if err != nil {
				// Invalid token, continue without authentication
				return next(c)
			}

			// Store user info in context
			c.Set(UserIDCtxKey, claims.UserID.String())
			c.Set(UserEmailCtxKey, claims.Email)
			c.Set(ClaimsCtxKey, claims)

			return next(c)
		}
	}
}

// CSRFMiddleware validates CSRF tokens for state-changing operations
func CSRFMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Only validate for state-changing operations
			switch c.Request().Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				csrfToken := c.Request().Header.Get("X-CSRF-Token")
				if csrfToken == "" {
					return c.JSON(http.StatusForbidden, map[string]string{
						"error": "missing csrf token",
					})
				}
				c.Set("csrf_token", csrfToken)
			}

			return next(c)
		}
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID := c.Get(UserIDCtxKey)
	if userID == nil {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusUnauthorized, "user not found in context")
	}

	id, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	return id, nil
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(c echo.Context) string {
	email := c.Get(UserEmailCtxKey)
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetClaimsFromContext extracts token claims from context
func GetClaimsFromContext(c echo.Context) *token.TokenClaims {
	claims := c.Get(ClaimsCtxKey)
	if claims == nil {
		return nil
	}
	return claims.(*token.TokenClaims)
}
