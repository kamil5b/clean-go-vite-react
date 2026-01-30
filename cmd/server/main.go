package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kamil5b/clean-go-vite-react/backend/di"
	"github.com/kamil5b/clean-go-vite-react/backend/platform"
	"github.com/kamil5b/clean-go-vite-react/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg := platform.NewConfig()

	// Create dependency container
	container := di.NewContainer(cfg)
	e := container.Echo

	// Setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup CORS middleware if needed
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Register frontend handlers (dev proxy or static assets)
	web.RegisterHandlers(e)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		e.Logger.Info("Starting server on " + addr)
		if err := e.Start(addr); err != nil && err.Error() != "http: Server closed" {
			e.Logger.Fatal(err)
		}
	}()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	e.Logger.Info("Server shutdown complete")
}
