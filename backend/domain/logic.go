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

	// User
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Item
	CreateItem(ctx context.Context, item *Item) error
	GetItems(ctx context.Context, userID uuid.UUID) ([]*Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error)
	UpdateItem(ctx context.Context, item *Item) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
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

// ItemInfo represents item information for API responses
type ItemInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID uuid.UUID `json:"sub"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	jwt.StandardClaims
}

// Private methods

// Item operations

// CreateItem creates a new item
func (l *Logic) CreateItem(ctx context.Context, title, description string, userID uuid.UUID) (*Item, error) {
	item := &Item{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	if err := l.db.CreateItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetItems retrieves all items for a user
func (l *Logic) GetItems(ctx context.Context, userID uuid.UUID) ([]*Item, error) {
	return l.db.GetItems(ctx, userID)
}

// GetItemByID retrieves an item by ID
func (l *Logic) GetItemByID(ctx context.Context, id uuid.UUID) (*Item, error) {
	return l.db.GetItemByID(ctx, id)
}

// UpdateItem updates an existing item
func (l *Logic) UpdateItem(ctx context.Context, id uuid.UUID, title, description string) (*Item, error) {
	item, err := l.db.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}

	item.Title = title
	item.Description = description

	if err := l.db.UpdateItem(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem deletes an item
func (l *Logic) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return l.db.DeleteItem(ctx, id)
}

// ItemToInfo converts an Item to ItemInfo for API responses
func ItemToInfo(item *Item) *ItemInfo {
	return &ItemInfo{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		UserID:      item.UserID,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

// ItemsToInfoList converts a slice of Items to ItemInfo slice
func ItemsToInfoList(items []*Item) []*ItemInfo {
	result := make([]*ItemInfo, len(items))
	for i, item := range items {
		result[i] = ItemToInfo(item)
	}
	return result
}

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
