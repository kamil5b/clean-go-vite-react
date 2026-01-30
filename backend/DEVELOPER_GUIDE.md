# Internal Architecture Developer Guide

## Overview

This guide documents the architecture and structure of the `/internal` package, which follows **Clean Architecture** principles with clear separation of concerns across multiple layers.

## Architecture Layers

### 1. **Model Layer** (`/backend/model`)
Defines domain models and data structures.

**Files:**
- `message.go` - Contains core domain entities:
  - `Message` - Message entity with ID, content, and timestamp
  - `User` - User entity with profile information
  - `Email` - Email record with delivery status

**Usage:**
- Pure data structures with JSON tags for serialization
- No business logic, no dependencies
- Shared across all layers

---

### 2. **Repository Layer** (`/backend/repository`)
Defines data access abstractions and implements persistence.

**Files:**
- `interfaces.go` - Repository interfaces:
  - `MessageRepository` - Message data access operations
  - `EmailRepository` - Email log storage operations
  - `UserRepository` - User CRUD operations
- `memory.go` - In-memory implementations for repositories
- `memory_test.go` - Unit tests for in-memory implementations

**Key Concepts:**
- Interface-based design for easy testing and swapping implementations
- Context-aware methods for request cancellation
- Returns errors for failed operations

**Example Interface:**
```go
type MessageRepository interface {
    GetMessage(ctx context.Context) (string, error)
}
```

---

### 3. **Service Layer** (`/backend/service`)
Implements business logic and orchestrates repository operations.

**Files:**
- `message.go` - Message business logic
  - `MessageService` interface
  - `messageService` implementation
  - Operations: GetMessage
- `email.go` - Email handling service
  - `EmailService` interface
  - Placeholder for SendEmail implementation
- `health.go` - Health check service
- `*_test.go` - Service unit tests

**Key Concepts:**
- Interfaces for dependency injection and testing
- Context support for graceful shutdown
- Encapsulates business rules and validation
- Depends on repository layer, not directly on storage

**Example Service:**
```go
type MessageService interface {
    GetMessage(ctx context.Context) (string, error)
}

type messageService struct{}
```

---

### 4. **Handler/API Layer** (`/backend/api`)
Exposes HTTP endpoints and handles HTTP requests/responses.

**Structure:**
```
/backend/api/
├── router.go           # Route configuration
├── handler/
│   ├── message.go      # Message HTTP handler
│   ├── health.go       # Health check handler
│   └── *_test.go       # Handler unit tests
└── integration_test.go # Integration tests
```

**Router Setup:**
- Configures API routes with base path `/api`
- Registers handlers with services
- Supports nested route groups

**Handler Pattern:**
- Each handler receives service dependencies
- Decodes request, calls service, encodes response
- Handles HTTP status codes and error responses

**Example Handler:**
```go
type MessageHandler struct {
    service service.MessageService
}

func (h *MessageHandler) GetMessage(c echo.Context) error {
    // Implementation
}
```

---

### 5. **Web Layer** (`/backend/web`)
Serves static web assets and frontend files.

**Files:**
- `web.go` - Web server setup and static file serving

---

### 6. **Worker/Job Layer** (`/backend/worker`)
Implements background job processing using Asynq.

**Files:**
- `email_processor.go` - Email notification task processor
  - `EmailProcessor` - Processes email tasks from queue
  - `ProcessTask` - Handles individual email tasks
- `email_processor_test.go` - Unit tests
- `email_processor_e2e_test.go` - End-to-end tests

**Key Concepts:**
- Decouples long-running operations from HTTP requests
- Uses Asynq for reliable task queuing
- Payload marshaling/unmarshaling with JSON
- Retry logic with configurable delays

**Example Worker:**
```go
type EmailProcessor struct {
    service service.EmailService
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
    // Parse payload and process
}
```

---

### 7. **Task Layer** (`/backend/task`)
Defines background job payloads and task types.

**Files:**
- `tasks.go` - Task payload structures and factory functions
- `tasks_test.go` - Task tests

---

### 8. **Dependency Injection** (`/backend/di`)
Manages dependency initialization and wiring.

**Files:**
- `container.go` - Main DI container
  - `Container` - Holds all application dependencies
  - `Services` - Service layer dependencies
  - `Handlers` - Handler layer dependencies
  - `NewContainer()` - Factory for creating the container
- `container_test.go` - Container tests

**Initialization Flow:**
1. Create Echo instance
2. Initialize services
3. Initialize handlers
4. Setup routes
5. Return populated container

**Example:**
```go
container := di.NewContainer(cfg)
container.Echo.Start(":8080")
```

---

### 9. **Platform/Configuration** (`/backend/platform`)
Infrastructure configuration and external service setup.

**Files:**
- `config.go` - Configuration management
  - `Config` - Main configuration struct
  - `ServerConfig` - HTTP server settings
  - `DatabaseConfig` - Database connection settings
  - `RedisConfig` - Redis connection settings
  - `AsynqConfig` - Job queue settings
  - `NewConfig()` - Load from environment
- `asynq.go` - Asynq job queue integration
- `config_test.go` - Configuration tests

**Configuration Source:**
- Environment variables with sensible defaults
- Parsing utilities: `getEnv()`, `getEnvInt()`, `getEnvBool()`, `getEnvDuration()`

**Environment Variables:**
```
SERVER_PORT=8080
SERVER_HOST=localhost
DATABASE_DSN=postgres://...
REDIS_HOST=localhost
REDIS_PORT=6379
ASYNQ_ENABLED=true
ASYNQ_REDIS_ADDR=localhost:6379
ASYNQ_CONCURRENCY=10
ASYNQ_MAX_RETRIES=3
```

---

## Data Flow

### HTTP Request Flow
```
HTTP Request
    ↓
Router (/backend/api/router.go)
    ↓
Handler (/backend/api/handler/*.go)
    ↓
Service (/backend/service/*.go)
    ↓
Repository (/backend/repository/*.go)
    ↓
HTTP Response
```

### Background Job Flow
```
Task Enqueue
    ↓
Task Queue (Redis via Asynq)
    ↓
Worker (/backend/worker/*.go)
    ↓
Service (/backend/service/*.go)
    ↓
Repository (/backend/repository/*.go)
    ↓
Completion (Success/Retry/Failed)
```

---

## Testing Strategy

### Unit Tests
- Located alongside source files with `_test.go` suffix
- Test individual components in isolation
- Use mocks/interfaces for dependencies

**Test Files:**
- `repository/*_test.go` - Repository tests
- `service/*_test.go` - Service tests
- `api/handler/*_test.go` - Handler tests
- `di/container_test.go` - Container tests
- `platform/config_test.go` - Configuration tests
- `task/tasks_test.go` - Task tests
- `worker/*_test.go` - Worker tests

### Integration Tests
- `api/integration_test.go` - Full API flow tests
- `worker/*_e2e_test.go` - End-to-end worker tests

---

## Adding a New Feature

### 1. Create Domain Models
Add to `/backend/model/` if new domain concepts needed.

### 2. Define Repository Interface
Add interface to `/backend/repository/interfaces.go`.

### 3. Implement Repository
Add implementation to `/backend/repository/memory.go` (or new file).

### 4. Create Service Interface and Implementation
Add interface and struct to `/backend/service/feature.go`.

### 5. Create Handler
Add HTTP handler to `/backend/api/handler/feature.go`.

### 6. Register Route
Update `/backend/api/router.go` to register the new endpoint.

### 7. Wire Dependencies
Update `/backend/di/container.go` to inject dependencies.

### 8. Write Tests
Add unit tests and integration tests for new components.

---

## Dependency Direction (Clean Architecture)

```
api/handler ↓
  ↓
service ↓
  ↓
repository ↓
  ↓
model
(no dependencies)
```

- **Model Layer**: No dependencies
- **Repository Layer**: Depends only on model
- **Service Layer**: Depends on model and repository
- **Handler/API Layer**: Depends on service
- **DI Container**: Wires everything together
- **Platform**: Configuration only, no business logic

---

## Best Practices

### 1. **Always Use Interfaces**
Define interfaces for services and repositories to enable testing and loose coupling.

### 2. **Context Propagation**
Pass `context.Context` through all layers for request cancellation and timeouts.

### 3. **Error Handling**
Return errors explicitly; don't panic in production code.

### 4. **Dependency Injection**
All dependencies should be injected via constructor functions or struct fields.

### 5. **Single Responsibility**
Each struct should have one reason to change.

### 6. **No Circular Dependencies**
Dependencies should flow in one direction (toward model layer).

### 7. **Test Coverage**
Aim for >80% test coverage for critical paths.

### 8. **Configuration Externalization**
Use environment variables for all configuration values.

---

## Common Patterns

### Service with Multiple Operations
```go
type UserService interface {
    Create(ctx context.Context, user *model.User) (string, error)
    GetByID(ctx context.Context, id string) (*model.User, error)
    Update(ctx context.Context, id string, user *model.User) error
    Delete(ctx context.Context, id string) error
}
```

### Handler Error Response
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func (h *Handler) handleError(c echo.Context, statusCode int, message string) error {
    return c.JSON(statusCode, ErrorResponse{
        Error:   "request_failed",
        Message: message,
    })
}
```

### Service with Repository Dependency
```go
type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
    return s.repo.FindByID(ctx, id)
}
```

---

## Package Imports Convention

```go
// Typical import structure in handlers
import (
    "context"
    "net/http"
    
    "github.com/labstack/echo/v4"
    
    "github.com/kamil5b/clean-go-vite-react/backend/model"
    "github.com/kamil5b/clean-go-vite-react/backend/service"
)
```

---

## Environment Configuration

Create `.env` file in project root:
```bash
SERVER_PORT=8080
SERVER_HOST=localhost
DATABASE_DSN=postgres://user:pass@localhost/dbname
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
REDIS_PASSWORD=
ASYNQ_ENABLED=true
ASYNQ_REDIS_ADDR=localhost:6379
ASYNQ_CONCURRENCY=10
ASYNQ_MAX_RETRIES=3
```

---

## Running Tests

```bash
# Run all tests
go test ./backend/...

# Run tests with coverage
go test -cover ./backend/...

# Run specific package tests
go test ./backend/service/...

# Run tests with verbose output
go test -v ./backend/...
```

---

## External Dependencies

- **Echo**: HTTP web framework
- **Asynq**: Background job queue with Redis
- **Database**: Configurable via DSN (PostgreSQL supported)
- **Redis**: For caching and job queue

---
