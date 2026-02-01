package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kamil5b/clean-go-vite-react/backend/api"
	"github.com/kamil5b/clean-go-vite-react/backend/domain"
	"github.com/kamil5b/clean-go-vite-react/backend/infra"
	web "github.com/kamil5b/clean-go-vite-react/embedder"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Initialize database
	dbConfig := infra.Config{
		Type:            getEnv("DB_TYPE", "sqlite"),
		DSN:             getEnv("DB_DSN", "dev.db"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := infra.NewDB(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}

	// Initialize domain logic
	logic := domain.NewLogic(db)

	// Initialize Echo
	e := echo.New()

	// Setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	// Register API routes
	api.SetupRoutes(e, logic)

	// Register frontend handler (dev proxy or static assets)
	web.RegisterHandlers(e)

	// Start server in a goroutine
	go func() {
		host := getEnv("SERVER_HOST", "0.0.0.0")
		port := getEnv("SERVER_PORT", "8080")
		addr := fmt.Sprintf("%s:%s", host, port)
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
