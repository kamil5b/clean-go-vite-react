package handler

import (
	"net/http"

	"github.com/kamil5b/clean-go-vite-react/backend/api/middleware"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/service/csrf"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
	userSvc "github.com/kamil5b/clean-go-vite-react/backend/service/user"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService  userSvc.UserService
	tokenService token.TokenService
	csrfService  csrf.CSRFService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService userSvc.UserService, tokenService token.TokenService, csrfService csrf.CSRFService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
		csrfService:  csrfService,
	}
}

// Register handles POST /api/auth/register requests
func (h *UserHandler) Register(c echo.Context) error {
	req := &request.RegisterUserRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email, password, and name are required",
		})
	}

	// Register user
	resp, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Set HTTP-only cookies
	c.SetCookie(&http.Cookie{
		Name:     middleware.AccessTokenCookie,
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	return c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/auth/login requests
func (h *UserHandler) Login(c echo.Context) error {
	req := &request.LoginRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email and password are required",
		})
	}

	// Login user
	resp, err := h.userService.Login(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// Generate refresh token
	refreshToken, err := h.tokenService.GenerateRefreshToken(resp.User.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate refresh token",
		})
	}

	// Set HTTP-only cookies for both access and refresh tokens
	c.SetCookie(&http.Cookie{
		Name:     middleware.AccessTokenCookie,
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	return c.JSON(http.StatusOK, resp)
}

// Refresh handles POST /api/auth/refresh requests
func (h *UserHandler) Refresh(c echo.Context) error {
	// Get refresh token from cookie
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing refresh token",
		})
	}

	// Refresh access token
	resp, err := h.userService.Refresh(c.Request().Context(), cookie.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// Set new access token cookie
	c.SetCookie(&http.Cookie{
		Name:     middleware.AccessTokenCookie,
		Value:    resp.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	return c.JSON(http.StatusOK, resp)
}

// Logout handles POST /api/auth/logout requests
func (h *UserHandler) Logout(c echo.Context) error {
	// Clear cookies
	c.SetCookie(&http.Cookie{
		Name:     middleware.AccessTokenCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

// GetMe handles GET /api/auth/me requests (protected)
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get(middleware.UserIDCtxKey)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	user, err := h.userService.GetUser(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetCSRFToken handles GET /api/csrf requests
func (h *UserHandler) GetCSRFToken(c echo.Context) error {
	token, err := h.csrfService.GenerateToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate csrf token",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
