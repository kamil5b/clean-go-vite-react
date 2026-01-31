# Model Package Developer Guide

## Overview

The `model` package is the central location for all data structures used in the backend application. It organizes types into three main categories:

- **Entity**: Database models (GORM) representing actual database tables
- **Request**: Incoming API request payload structures
- **Response**: Outgoing API response payload structures

## Directory Structure

```
model/
├── entity/         # Database entities (GORM models)
├── request/        # API request DTOs
├── response/       # API response DTOs
└── README.md
```

## Adding New Models

### 1. Entity Models

Create a new file in `entity/` directory for your database model.

**Conventions:**
- Suffix the struct name with `Entity` (e.g., `UserEntity`)
- Include GORM tags for primary keys and soft deletes
- Add timestamp fields: `CreatedAt`, `UpdatedAt`, `DeletedAt`
- Use `uuid.UUID` for primary keys

**Example:**
```go
package entity

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 2. Request Models

Create a new file in `request/` directory for your API request payloads.

**Conventions:**
- Suffix the struct name with `Request` (e.g., `CreateProductRequest`)
- Include JSON tags for request body binding
- Only include fields that are expected from the client
- Use validation tags if using a validator (e.g., `validate:required`)

**Example:**
```go
package request

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
```

### 3. Response Models

Create a new file in `response/` directory for your API response payloads.

**Conventions:**
- Name the struct based on the action (e.g., `GetProduct`, `ListProducts`)
- Include JSON tags for response serialization
- Only include fields that should be exposed to the client
- Exclude sensitive data (passwords, secrets, etc.)

**Example:**
```go
package response

import "github.com/google/uuid"

type GetProduct struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}
```

## Common Patterns

### Response Wrappers

Use `CommonIDResponse` for simple ID-only responses:

```go
package response

type CommonIDResponse struct {
	ID uuid.UUID `json:"value"`
}
```

### Data Flow

Typical request → response flow:

1. **Client sends request** → `request.CreateProductRequest`
2. **Handler receives request** → Maps to `entity.ProductEntity` or business logic
3. **Handler returns response** → `response.GetProduct` with relevant fields

### Mapping Between Layers

Always map between Request → Entity → Response to:
- Prevent exposing unnecessary fields
- Maintain separation of concerns
- Control what clients can see/modify

**Example:**
```go
// request.CreateProductRequest → entity.ProductEntity → response.GetProduct
product := &entity.ProductEntity{
	ID:    uuid.New(),
	Name:  req.Name,
	Price: req.Price,
}

// Later, when responding:
return &response.GetProduct{
	ID:    product.ID,
	Name:  product.Name,
	Price: product.Price,
}
```

## Existing Models

### User
- **Entity**: `UserEntity` - Core user data with password storage
- **Request**: `RegisterUserRequest`, `LoginRequest` - User authentication payloads
- **Response**: `GetUser` - Safe user data (no password)

### Counter
- **Entity**: `CounterEntity` - Simple counter storage
- **Response**: `GetCounter` - Counter value for API responses

### Message
- **Entity**: `MessageEntity` - Key-value message storage
- **Response**: Message response structures

### Email
- **Entity**: `EmailEntity` - Email records
- **Request**: Email-related requests
- **Response**: Email response data

## Best Practices

1. **Naming**: Use consistent suffixes (`Entity`, `Request`, `Response`)
2. **Separation**: Never mix layers (request ≠ response ≠ entity)
3. **Security**: Exclude sensitive fields from responses
4. **Timestamps**: Always include `CreatedAt`, `UpdatedAt` for entities
5. **Soft Deletes**: Include `DeletedAt` field for auditing capability
6. **Type Safety**: Use `uuid.UUID` for IDs instead of strings
7. **Tags**: Always include JSON tags for serialization
8. **Documentation**: Add brief comments for public types

## Using Models in Handlers

```go
package handler

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

func CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Handle error
	}

	// Create entity
	product := &entity.ProductEntity{
		Name:  req.Name,
		Price: req.Price,
	}

	// Save to database (using service/repository)
	// ...

	// Map to response
	c.JSON(200, response.GetProduct{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	})
}
```

## Import Paths

```go
import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)
```

## Migration & Database Considerations

When adding a new `Entity` type:

1. Create the struct in `entity/` with proper GORM tags
2. Register it with GORM's `AutoMigrate()`
3. Ensure soft delete migrations if using `DeletedAt`
4. Test migrations in development environment

## Validation

For request validation, models should work with standard Go validation libraries. Add validation tags to request structs:

```go
type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required,min=1,max=255"`
	Price float64 `json:"price" validate:"required,gt=0"`
}
```
