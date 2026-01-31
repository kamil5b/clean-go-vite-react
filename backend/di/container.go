package di

import (
	"github.com/kamil5b/clean-go-vite-react/backend/api"
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/kamil5b/clean-go-vite-react/backend/platform"

	counterRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/counter"
	messageRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/message"

	counterSvc "github.com/kamil5b/clean-go-vite-react/backend/service/counter"
	emailSvc "github.com/kamil5b/clean-go-vite-react/backend/service/email"
	healthSvc "github.com/kamil5b/clean-go-vite-react/backend/service/health"
	messageSvc "github.com/kamil5b/clean-go-vite-react/backend/service/message"

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
	Email   emailSvc.EmailService
	Health  healthSvc.HealthService
	Counter counterSvc.CounterService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Message *handler.MessageHandler
	Health  *handler.HealthHandler
	Counter *handler.CounterHandler
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *platform.Config) *Container {
	// Initialize Echo
	e := echo.New()

	// Initialize repositories
	db := cfg.Database.Gorm
	if db == nil {
		db = platform.InitializeDatabase(cfg)
		cfg.Database.Gorm = db
	}

	counterRepository, _ := counterRepo.NewGORMCounterRepository(db)
	messageRepository, _ := messageRepo.NewGORMMessageRepository(db)

	services := &Services{
		Message: messageSvc.NewMessageService(messageRepository),
		Email:   emailSvc.NewEmailService(),
		Health:  healthSvc.NewHealthService(),
		Counter: counterSvc.NewCounterService(counterRepository),
	}

	// Initialize handlers
	handlers := &Handlers{
		Message: handler.NewMessageHandler(services.Message),
		Health:  handler.NewHealthHandler(services.Health),
		Counter: handler.NewCounterHandler(services.Counter),
	}

	// Setup routes with dependencies
	api.SetupRoutes(e, *handlers.Message, *handlers.Counter)
	e.GET("/api/health", handlers.Health.Check)

	return &Container{
		Config:   cfg,
		Echo:     e,
		Services: services,
	}
}
