package di

import (
	"github.com/kamil5b/clean-go-vite-react/internal/api"
	"github.com/kamil5b/clean-go-vite-react/internal/platform"
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
}

// NewContainer creates and initializes a new dependency container
func NewContainer(cfg *platform.Config) *Container {
	// Initialize Echo
	e := echo.New()

	// Initialize services
	services := &Services{
		Message: service.NewMessageService(),
	}

	// Setup routes with dependencies
	api.SetupRoutes(e, services.Message)

	return &Container{
		Config:   cfg,
		Echo:     e,
		Services: services,
	}
}
