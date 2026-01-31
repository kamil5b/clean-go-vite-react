package di

import (
	"os"
	"time"

	"github.com/kamil5b/clean-go-vite-react/backend/api"
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/kamil5b/clean-go-vite-react/backend/platform"

	counterRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/counter"
	messageRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/message"
	userRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/user"

	counterSvc "github.com/kamil5b/clean-go-vite-react/backend/service/counter"
	csrfSvc "github.com/kamil5b/clean-go-vite-react/backend/service/csrf"
	healthSvc "github.com/kamil5b/clean-go-vite-react/backend/service/health"
	messageSvc "github.com/kamil5b/clean-go-vite-react/backend/service/message"
	tokenSvc "github.com/kamil5b/clean-go-vite-react/backend/service/token"
	userSvc "github.com/kamil5b/clean-go-vite-react/backend/service/user"

	"github.com/labstack/echo/v4"
)

// Container holds all application dependencies
type Container struct {
	Config   *platform.Config
	Echo     *echo.Echo
	Services *Services
}

// Services holds all service layer dependencies
type Services struct {
	Message messageSvc.MessageService
	Health  healthSvc.HealthService
	Counter counterSvc.CounterService
	User    userSvc.UserService
	Token   tokenSvc.TokenService
	CSRF    csrfSvc.CSRFService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Message *handler.MessageHandler
	Health  *handler.HealthHandler
	Counter *handler.CounterHandler
	User    *handler.UserHandler
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *platform.Config) *Container {
	// Initialize Echo
	e := echo.New()

	// Initialize database
	db := cfg.Database.Gorm
	if db == nil {
		db = platform.InitializeDatabase(cfg)
		cfg.Database.Gorm = db
	}

	// Initialize repositories
	counterRepository, _ := counterRepo.NewGORMCounterRepository(db)
	messageRepository, _ := messageRepo.NewGORMMessageRepository(db)
	userRepository, _ := userRepo.NewGORMUserRepository(db)

	// Initialize token service with configuration from environment
	tokenConfig := tokenSvc.TokenConfig{
		AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret-key-change-in-production"),
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret-key-change-in-production"),
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	tokenService := tokenSvc.NewTokenService(tokenConfig)

	// Initialize services
	services := &Services{
		Message: messageSvc.NewMessageService(messageRepository),
		Health:  healthSvc.NewHealthService(),
		Counter: counterSvc.NewCounterService(counterRepository),
		User:    userSvc.NewUserService(userRepository, tokenService),
		Token:   tokenService,
		CSRF:    csrfSvc.NewCSRFService(),
	}

	// Initialize handlers
	handlers := &Handlers{
		Message: handler.NewMessageHandler(services.Message),
		Health:  handler.NewHealthHandler(services.Health),
		Counter: handler.NewCounterHandler(services.Counter),
		User:    handler.NewUserHandler(services.User, services.Token, services.CSRF),
	}

	// Setup routes with dependencies
	api.SetupRoutes(e, *handlers.Message, *handlers.Counter, handlers.User, services.Token)
	e.GET("/api/health", handlers.Health.Check)

	return &Container{
		Config:   cfg,
		Echo:     e,
		Services: services,
	}
}

// Helper function to get environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
