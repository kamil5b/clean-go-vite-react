package di

import (
	"log"
	"os"
	"time"

	"github.com/kamil5b/clean-go-vite-react/backend/api"
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/kamil5b/clean-go-vite-react/backend/platform"

	counterRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/counter"
	invoiceRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/invoice"
	itemRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/item"
	messageRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/message"
	tagRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/tag"
	userRepo "github.com/kamil5b/clean-go-vite-react/backend/repository/implementations/user"

	counterSvc "github.com/kamil5b/clean-go-vite-react/backend/service/counter"
	csrfSvc "github.com/kamil5b/clean-go-vite-react/backend/service/csrf"
	healthSvc "github.com/kamil5b/clean-go-vite-react/backend/service/health"
	invoiceSvc "github.com/kamil5b/clean-go-vite-react/backend/service/invoice"
	itemSvc "github.com/kamil5b/clean-go-vite-react/backend/service/item"
	messageSvc "github.com/kamil5b/clean-go-vite-react/backend/service/message"
	tagSvc "github.com/kamil5b/clean-go-vite-react/backend/service/tag"
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
	Item    itemSvc.ItemService
	Tag     tagSvc.TagService
	Invoice invoiceSvc.InvoiceService
}

// Handlers holds all HTTP handler dependencies
type Handlers struct {
	Message *handler.MessageHandler
	Health  *handler.HealthHandler
	Counter *handler.CounterHandler
	User    *handler.UserHandler
	Item    *handler.ItemHandler
	Tag     *handler.TagHandler
	Invoice *handler.InvoiceHandler
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
	counterRepository, err := counterRepo.NewGORMCounterRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize counter repository: %v", err)
	}

	messageRepository, err := messageRepo.NewGORMMessageRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize message repository: %v", err)
	}

	userRepository, err := userRepo.NewGORMUserRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}

	itemRepository, err := itemRepo.NewGORMItemRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize item repository: %v", err)
	}

	tagRepository, err := tagRepo.NewGORMTagRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize tag repository: %v", err)
	}

	invoiceRepository, err := invoiceRepo.NewGORMInvoiceRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize invoice repository: %v", err)
	}

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
		Item:    itemSvc.NewItemService(itemRepository),
		Tag:     tagSvc.NewTagService(tagRepository),
		Invoice: invoiceSvc.NewInvoiceService(invoiceRepository, tagRepository),
	}

	// Initialize handlers
	handlers := &Handlers{
		Message: handler.NewMessageHandler(services.Message),
		Health:  handler.NewHealthHandler(services.Health),
		Counter: handler.NewCounterHandler(services.Counter),
		User:    handler.NewUserHandler(services.User, services.Token, services.CSRF),
		Item:    handler.NewItemHandler(services.Item),
		Tag:     handler.NewTagHandler(services.Tag),
		Invoice: handler.NewInvoiceHandler(services.Invoice),
	}

	// Setup routes with dependencies
	api.SetupRoutes(e, *handlers.Message, *handlers.Counter, handlers.User, services.Token, handler.NewNotFoundHandler(), handlers.Item, handlers.Tag, handlers.Invoice)
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
