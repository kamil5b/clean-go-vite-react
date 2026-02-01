package domain

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Logic holds business logic and database access
type Logic struct {
	db        DB
	jwtSecret string
	jwtExpiry time.Duration
}

// DB interface for data access (infra layer)
type DB interface {
	// Message
	GetMessage(ctx context.Context, key string) (*Message, error)

	// Counter
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)

	// User
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// NewLogic creates a new domain logic instance
func NewLogic(db DB) *Logic {
	return &Logic{
		db:        db,
		jwtSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		jwtExpiry: 15 * time.Minute,
	}
}

// Message operations

// GetMessage retrieves the default message
func (l *Logic) GetMessage(ctx context.Context) (string, error) {
	msg, err := l.db.GetMessage(ctx, "default")
	if err != nil {
		return "", err
	}
	return msg.Value, nil
}

// Counter operations

// GetCounter returns the current counter value
func (l *Logic) GetCounter(ctx context.Context) (int, error) {
	return l.db.GetCounter(ctx)
}

// IncrementCounter increments and returns the new counter value
func (l *Logic) IncrementCounter(ctx context.Context) (int, error) {
	return l.db.IncrementCounter(ctx)
}

// User operations

// RegisterUser creates a new user account
func (l *Logic) RegisterUser(ctx context.Context, email, password, name string) (*UserInfo, string, error) {
	// Check if user exists
	existing, err := l.db.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	if err := l.db.CreateUser(ctx, user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := l.generateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, "", err
	}

	userInfo := &UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return userInfo, token, nil
}

// LoginUser authenticates a user
func (l *Logic) LoginUser(ctx context.Context, email, password string) (*UserInfo, string, error) {
	// Find user
	user, err := l.db.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate token
	token, err := l.generateToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, "", err
	}

	userInfo := &UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return userInfo, token, nil
}

// GetUserByID retrieves a user by ID
func (l *Logic) GetUserByID(ctx context.Context, id uuid.UUID) (*UserInfo, error) {
	user, err := l.db.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (l *Logic) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(l.jwtSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return uuid.Nil, errors.New("invalid token")
}

// Helper types

// UserInfo represents user information without sensitive data
type UserInfo struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	jwt.StandardClaims
}

// Private methods

func (l *Logic) generateToken(userID uuid.UUID, email, name string) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(l.jwtExpiry).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.jwtSecret))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
