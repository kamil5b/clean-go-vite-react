# Migration Plan

This document provides a step-by-step guide for migrating this forked repository (Go + Vite + React) to align with the full architecture defined in `ARCHITECTURE.md`. The migration builds upon the existing clean separation between frontend and backend, adding runtime boundaries, service layers, and async processing.

---

## Current State (From Forked Repository)

**What we already have:**
- Go (Echo) backend with `/api/*` routes
- Vite + React frontend with HMR in development
- `internal/web` layer for optional frontend hosting
- `internal/api` for API routes & handlers
- Single-binary deployment via `go:embed`
- Clean separation: backend doesn't know React exists
- Dev mode: Vite dev server + Go server
- Prod mode: Embedded assets in Go binary

**What we need to add:**
- Separated worker runtime for async processing
- Service layer (transport-agnostic business logic)
- Repository layer (data access abstraction)
- Task contracts and Asynq integration
- Platform layer for infrastructure
- Model layer for domain objects
- Enhanced dependency injection

---

## Migration Overview

**Goal**: Enhance the current clean architecture with:
- Separated server and worker runtimes
- Transport-agnostic service layer
- Async processing with Asynq
- Maintained frontend detachability
- Strict dependency rules
- Shared business logic across runtimes

**Approach**: Incremental migration that preserves existing functionality and developer experience.

---

## Phase 0: Assessment & Preparation

### 0.1 Audit Current Structure
**Current state from repository:**
- [ ] `cmd/server/main.go` - Application entrypoint exists
- [ ] `internal/api/` - API routes & handlers exist
- [ ] `internal/web/` - Optional web hosting layer exists (dev_proxy.go, static.go, embedded_assets.go)
- [ ] `frontend/` - Standalone Vite + React app exists
- [ ] Document all existing HTTP handlers in `internal/api/`
- [ ] Identify business logic currently in handlers
- [ ] List all database operations (if any)
- [ ] Identify any operations that should be async
- [ ] Verify frontend-backend separation (should be good already)

### 0.2 Setup Infrastructure
- [ ] Add Redis to development environment (docker-compose or local)
- [ ] Install Asynq: `go get github.com/hibiken/asynq`
- [ ] Review existing `Makefile` and add new commands for worker
- [ ] Review existing `Dockerfile` for multi-runtime support
- [ ] Create migration branch
- [ ] Ensure `make dev` still works for current setup

### 0.3 Documentation
- [ ] Create `DETACH.md` placeholder
- [ ] Document breaking changes plan
- [ ] Backup current working state

**Milestone**: Clear understanding of current state and migration scope

---

## Phase 1: Directory Restructuring

### 1.1 Create New Directory Structure

**Starting point (current):**
```
├── cmd/
│   └── server/
│       └── main.go          ✓ EXISTS
├── internal/
│   ├── api/                 ✓ EXISTS (routes & handlers)
│   └── web/                 ✓ EXISTS (dev_proxy.go, static.go, embedded_assets.go)
└── frontend/                ✓ EXISTS (Vite + React)
```

**Target structure (add these):**
```
├── cmd/
│   ├── server/              ✓ EXISTS
│   │   └── main.go
│   └── worker/              ← ADD
│       └── main.go
├── internal/
│   ├── api/                 ✓ EXISTS - REFACTOR
│   │   ├── handler/         ← ORGANIZE
│   │   └── router.go        ← ADD
│   ├── service/             ← ADD
│   ├── repository/          ← ADD
│   │   └── interfaces/
│   ├── task/                ← ADD
│   ├── worker/              ← ADD
│   ├── model/               ← ADD
│   ├── platform/            ← ADD
│   └── web/                 ✓ EXISTS - KEEP AS-IS
└── frontend/                ✓ EXISTS - KEEP AS-IS
```

### 1.2 Organize Existing Files
- [ ] Create all new directories (service, repository, task, worker, model, platform)
- [ ] **Keep `internal/web/` unchanged** (dev_proxy.go, static.go, embedded_assets.go)
- [ ] Organize existing handlers in `internal/api/` into `internal/api/handler/`
- [ ] Extract any existing models to `internal/model/`
- [ ] Frontend already in place at `frontend/` - no changes needed
- [ ] Keep old handler structure temporarily for reference

### 1.3 Update Import Paths
- [ ] Search and replace old import paths
- [ ] Fix compilation errors
- [ ] Verify tests still pass

**Milestone**: New structure in place, existing code compiles

---

## Phase 2: Extract Service Layer

### 2.1 Identify Business Logic
- [ ] Review all handlers for business rules
- [ ] List operations that modify state
- [ ] Identify validation logic
- [ ] Map out data transformations

### 2.2 Create Service Interfaces
```go
// Example structure
package service

type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id string) (*User, error)
    UpdateUser(ctx context.Context, id string, req UpdateUserRequest) error
    DeleteUser(ctx context.Context, id string) error
}
```

### 2.3 Implement Services
- [ ] Create service structs in `internal/service/`
- [ ] Move business logic from handlers to services
- [ ] Services should accept primitives/domain objects
- [ ] Services should NOT import Echo, Asynq, or HTTP types
- [ ] Add comprehensive unit tests

### 2.4 Update Handlers
- [ ] Handlers become thin wrappers
- [ ] Parse request → call service → serialize response
- [ ] Remove business logic from handlers
- [ ] Keep validation at handler level

**Milestone**: Business logic extracted, handlers are thin, services are testable

---

## Phase 3: Extract Repository Layer

### 3.1 Define Repository Interfaces
```go
// Example structure
package repository

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    FindByID(ctx context.Context, id string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id string) error
}
```

### 3.2 Implement Repositories
- [ ] Create repository implementations in `internal/repository/postgres/` (or appropriate DB)
- [ ] Move all SQL/database logic from services to repositories
- [ ] Repositories should be pure data access
- [ ] Add integration tests

### 3.3 Wire Repositories to Services
- [ ] Services depend on repository interfaces
- [ ] Inject repositories via constructors
- [ ] Remove direct database calls from services

**Milestone**: Data access abstracted, services are database-agnostic

---

## Phase 4: Setup Platform Layer

### 4.1 Create Platform Package
- [ ] Create `internal/platform/config.go` for configuration
- [ ] Create `internal/platform/database.go` for DB connections
- [ ] Create `internal/platform/redis.go` for Redis connections
- [ ] Create `internal/platform/asynq.go` for Asynq setup

### 4.2 Configuration Management
```go
// Example structure
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Asynq    AsynqConfig
}
```

### 4.3 Connection Pooling
- [ ] Setup database connection pooling
- [ ] Setup Redis connection pooling
- [ ] Add health check endpoints
- [ ] Add graceful shutdown handlers

**Milestone**: Infrastructure properly abstracted and configurable

---

## Phase 5: Implement Async Processing

### 5.1 Setup Redis
- [ ] Add Redis to docker-compose or deployment
- [ ] Configure Redis connection settings
- [ ] Test Redis connectivity

### 5.2 Define Task Contracts
```go
// internal/task/tasks.go
package task

const (
    TypeEmailNotification = "email:notification"
    TypeDataExport        = "data:export"
)

type EmailNotificationPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}
```

### 5.3 Create Task Enqueuing (Server Side)
- [ ] Create Asynq client in platform layer
- [ ] Add methods to enqueue tasks
- [ ] Update handlers to enqueue async work
- [ ] Remove synchronous execution of long-running tasks

### 5.4 Create Task Processors (Worker Side)
```go
// internal/worker/processor.go
package worker

type EmailProcessor struct {
    service service.EmailService
}

func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
    var payload task.EmailNotificationPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    return p.service.SendEmail(ctx, payload.To, payload.Subject, payload.Body)
}
```

### 5.5 Create Worker Runtime
- [ ] Create `cmd/worker/main.go`
- [ ] Initialize Asynq server
- [ ] Register task processors
- [ ] Setup retry policies and error handling
- [ ] Add worker monitoring

**Milestone**: Async processing functional, long tasks moved out of request cycle

---

## Phase 6: Split Server Runtime

### 6.1 Refactor Server Main
- [ ] Move server setup to `cmd/server/main.go`
- [ ] Initialize dependencies (DB, Redis, services)
- [ ] Setup Echo router
- [ ] Register API routes
- [ ] Setup middleware (CORS, logging, etc.)

### 6.2 API Router Organization
```go
// internal/api/router.go
package api

func SetupRoutes(e *echo.Echo, deps *Dependencies) {
    api := e.Group("/api")
    
    // Health check
    api.GET("/health", deps.HealthHandler.Check)
    
    // Feature routes
    users := api.Group("/users")
    users.POST("", deps.UserHandler.Create)
    users.GET("/:id", deps.UserHandler.Get)
}
```

### 6.3 Dependency Injection
- [ ] Create dependency container
- [ ] Wire all components together
- [ ] Use constructor injection
- [ ] Avoid global state

**Milestone**: Server runtime cleanly separated and properly wired

---

## Phase 7: Frontend Integration

### 7.1 Frontend API Client
- [ ] Review existing frontend API communication
- [ ] Create or enhance `frontend/lib/api.ts` for HTTP client
- [ ] Define API endpoints as constants
- [ ] Add error handling
- [ ] Add request/response interceptors

### 7.2 Environment Configuration
**Already configured in vite.config.ts:**
```typescript
// Existing proxy configuration
server: {
  proxy: {
    '/api': 'http://localhost:3000'
  }
}
```
- [ ] Verify API_BASE_URL configuration
- [ ] Ensure `/api/*` prefix is consistent

### 7.3 Development Setup (Already Working)
**Current setup to maintain:**
- `internal/web/dev_proxy.go` ✓ EXISTS - forwards to Vite dev server
- `internal/web/web.go` ✓ EXISTS - dev vs prod selection
- `make dev` runs both Go and Vite servers
- [ ] Verify HMR still works after adding services
- [ ] Test dev mode with new architecture

### 7.4 Production Build (Already Working)
**Current setup to maintain:**
- `internal/web/embedded_assets.go` ✓ EXISTS - `go:embed` for dist
- `internal/web/static.go` ✓ EXISTS - serves embedded files
- `make build` compiles single binary
- [ ] Verify production build works with new architecture
- [ ] Test embedded assets serving
- [ ] Verify Docker build still works

**Milestone**: Frontend properly integrated with clear API boundaries

---

## Phase 8: Testing & Validation

### 8.1 Unit Tests
- [ ] Services: Test business logic in isolation
- [ ] Repositories: Test with test database
- [ ] Handlers: Test with mocked services
- [ ] Workers: Test with mocked services

### 8.2 Integration Tests
- [ ] API endpoints: Full request/response cycle
- [ ] Async tasks: Enqueue and process
- [ ] Database operations: Real DB interactions

### 8.3 End-to-End Tests
- [ ] Critical user flows
- [ ] Async workflows
- [ ] Error scenarios

### 8.4 Performance Testing
- [ ] Load test API endpoints
- [ ] Monitor worker throughput
- [ ] Check database connection pooling
- [ ] Verify Redis performance

**Milestone**: Comprehensive test coverage, system validated

---

## Phase 9: Documentation & Cleanup

### 9.1 Update Documentation
- [ ] Update `README.md` with new architecture (worker runtime, async tasks)
- [ ] Keep existing README sections on dev mode, HMR, and embedding
- [ ] Document API endpoints (consider OpenAPI/Swagger)
- [ ] Document async task types and payloads
- [ ] Enhance existing `DETACH.md` (if exists) or create new one
- [ ] Document Makefile commands for server and worker
- [ ] Add inline code documentation

### 9.2 Code Cleanup
- [ ] Remove old/unused code
- [ ] Remove temporary migration code
- [ ] Fix linting issues
- [ ] Format code consistently

### 9.3 Developer Experience
- [ ] Update Makefile with new commands
- [ ] Create docker-compose for local development
- [ ] Add VS Code/IDE configurations
- [ ] Document common development workflows

**Milestone**: Clean, well-documented codebase

---

## Phase 10: Deployment

### 10.1 Deployment Preparation
- [ ] Update existing `Dockerfile` for server runtime (multi-stage already exists)
- [ ] Create new `Dockerfile.worker` for worker runtime
- [ ] Or enhance existing Dockerfile to support both runtimes
- [ ] Setup environment variables for both runtimes
- [ ] Configure logging and monitoring
- [ ] Setup health checks for server and worker

### 10.2 Infrastructure
- [ ] Deploy Redis instance
- [ ] Configure database
- [ ] Setup load balancer (if needed)
- [ ] Configure auto-scaling for workers

### 10.3 CI/CD Pipeline
- [ ] Build server binary
- [ ] Build worker binary
- [ ] Build frontend
- [ ] Run tests
- [ ] Deploy to staging
- [ ] Deploy to production

### 10.4 Monitoring
- [ ] Application metrics
- [ ] Asynq queue metrics
- [ ] Error tracking
- [ ] Performance monitoring

**Milestone**: System deployed and running in production

---

## Validation Checklist

After migration, verify:

### Architecture Compliance
- [ ] No service → handler dependencies
- [ ] No service → Asynq dependencies
- [ ] No repository → service dependencies
- [ ] Services are transport-agnostic
- [ ] Clear runtime boundaries

### Functionality
- [ ] All API endpoints working
- [ ] Async tasks processing correctly
- [ ] Frontend communicating properly
- [ ] Database operations successful
- [ ] Error handling working

### Detachment Readiness
- [ ] Frontend can be deployed separately
- [ ] Worker can scale independently
- [ ] New runtime can be added easily
- [ ] Async system can be replaced

### Quality
- [ ] All tests passing
- [ ] No critical linting issues
- [ ] Documentation complete
- [ ] Performance acceptable

---

## Rollback Plan

If issues arise during migration:

1. **Keep old code in branch**: Don't delete old implementation until validated
2. **Feature flags**: Use flags to switch between old/new implementations
3. **Database migrations**: Keep reversible, test rollback
4. **Monitoring**: Watch metrics closely after each phase
5. **Gradual rollout**: Deploy to staging first, then production gradually

---

## Common Pitfalls to Avoid

1. **Circular dependencies**: Always check dependency direction
2. **God services**: Keep services focused and cohesive
3. **Leaky abstractions**: Don't let transport concerns leak into services
4. **Premature optimization**: Get it working, then optimize
5. **Skipping tests**: Write tests as you go, not at the end
6. **Big bang migration**: Migrate incrementally, validate frequently

---

## Success Criteria

The migration is successful when:

1. All functionality works as before
2. Architecture rules are followed
3. Tests pass with good coverage
4. Code is cleaner and more maintainable
5. Frontend can be detached with minimal effort
6. Worker can scale independently
7. New features are easier to add
8. Documentation is complete

---

## Timeline Estimation

- **Phase 0-1**: 1-2 days (Assessment & structure)
- **Phase 2-3**: 3-5 days (Service & repository extraction)
- **Phase 4**: 1-2 days (Platform layer)
- **Phase 5**: 3-4 days (Async processing)
- **Phase 6**: 2-3 days (Server runtime)
- **Phase 7**: 2-3 days (Frontend integration)
- **Phase 8**: 3-5 days (Testing)
- **Phase 9**: 2-3 days (Documentation & cleanup)
- **Phase 10**: 2-3 days (Deployment)

**Total**: 3-4 weeks for a complete migration (adjust based on codebase size)

---

## Next Steps

1. Review this plan with the team
2. Start Phase 0: Assessment
3. Create GitHub issues for each phase
4. Begin migration incrementally
5. Review and adjust plan as needed

---

## Preserving What Works

**DO NOT BREAK these existing features:**
- ✅ Vite HMR in development mode
- ✅ Single-binary deployment
- ✅ `make dev` experience
- ✅ `make build` producing unified binary
- ✅ Docker multi-stage build
- ✅ Frontend detachability (already clean)
- ✅ Backend independence from React

**The migration adds capabilities without removing what works.**

---

## Questions & Support

If you encounter issues or need clarification:

1. Review `ARCHITECTURE.md` for design principles
2. Review `README.md` for current implementation details
3. Check dependency rules section
4. Test `make dev` after each phase
5. Verify HMR still works
6. Ask for code review early and often

Remember: The goal is not perfection, but consistent improvement toward clean architecture while preserving the excellent developer experience already in place.