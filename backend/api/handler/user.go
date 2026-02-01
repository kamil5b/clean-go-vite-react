package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	logic Logic
}

// NewUserHandler creates a new user handler
func NewUserHandler(logic Logic) *UserHandler {
	return &UserHandler{logic: logic}
}

// RegisterRequest represents registration input
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register
func (h *UserHandler) Register(c echo.Context) error {
	req := &RegisterRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email, password, and name are required",
		})
	}

	userInfo, token, err := h.logic.RegisterUser(c.Request().Context(), req.Email, req.Password, req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Set HTTP-only cookie
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":  userInfo,
		"token": token,
	})
}

// Login handles POST /api/auth/login
func (h *UserHandler) Login(c echo.Context) error {
	req := &LoginRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email and password are required",
		})
	}

	userInfo, token, err := h.logic.LoginUser(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// Set HTTP-only cookie
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  userInfo,
		"token": token,
	})
}

// Logout handles POST /api/auth/logout
func (h *UserHandler) Logout(c echo.Context) error {
	// Clear cookie
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged out successfully",
	})
}

// GetMe handles GET /api/auth/me (protected)
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid user id",
		})
	}

	user, err := h.logic.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
