package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kamil5b/clean-go-vite-react/backend/platform"
	"github.com/kamil5b/clean-go-vite-react/backend/service"
	"github.com/kamil5b/clean-go-vite-react/backend/task"
	"github.com/kamil5b/clean-go-vite-react/backend/worker"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg := platform.NewConfig()

	if !cfg.Asynq.Enabled {
		fmt.Println("Asynq is disabled in configuration")
		return
	}

	// Create Asynq server
	srv := platform.NewAsynqServer(&cfg.Asynq)

	// Create task multiplexer
	mux := platform.NewAsynqMux()

	// Initialize services
	emailService := service.NewEmailService()

	// Register task processors
	emailProcessor := worker.NewEmailProcessor(emailService)
	mux.HandleFunc(task.TypeEmailNotification, emailProcessor.ProcessTask)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		fmt.Println("Starting worker...")
		if err := srv.Run(mux); err != nil {
			fmt.Printf("Worker error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("Shutting down worker...")
	srv.Shutdown()

	fmt.Println("Worker shutdown complete")
}
