package di

import (
	"github.com/kamil5b/clean-go-vite-react/internal/api"
	"github.com/kamil5b/clean-go-vite-react/internal/api/handler"
	"github.com/kamil5b/clean-go-vite-react/internal/platform"
	"github.com/kamil5b/clean-go-vite-react/internal/repository"
	"github.com/kamil5b/clean-go-vite-react/internal/service"
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
	Message service.MessageService
	Email   service.EmailService
	Health  service.HealthService
	Counter service.CounterService
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
	counterRepo := repository.NewInMemoryCounterRepository()

	// Initialize services
	services := &Services{
		Message: service.NewMessageService(),
		Email:   service.NewEmailService(),
		Health:  service.NewHealthService(),
		Counter: service.NewCounterService(counterRepo),
	}

	// Initialize handlers
	handlers := &Handlers{
		Message: handler.NewMessageHandler(services.Message),
		Health:  handler.NewHealthHandler(services.Health),
		Counter: handler.NewCounterHandler(services.Counter),
	}

	// Setup routes with dependencies
	api.SetupRoutes(e, services.Message, services.Counter)
	e.GET("/api/health", handlers.Health.Check)

	return &Container{
		Config:   cfg,
		Echo:     e,
		Services: services,
	}
}
