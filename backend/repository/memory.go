package repository

import (
	"context"
	"fmt"
	"sync"
)

// InMemoryMessageRepository is an in-memory implementation of MessageRepository
type InMemoryMessageRepository struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewInMemoryMessageRepository creates a new in-memory message repository
func NewInMemoryMessageRepository() *InMemoryMessageRepository {
	return &InMemoryMessageRepository{
		data: make(map[string]string),
	}
}

// GetMessage returns a stored message
func (r *InMemoryMessageRepository) GetMessage(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if msg, ok := r.data["default"]; ok {
		return msg, nil
	}
	return "", fmt.Errorf("message not found")
}

// InMemoryEmailRepository is an in-memory implementation of EmailRepository
type InMemoryEmailRepository struct {
	mu      sync.RWMutex
	data    map[string]map[string]interface{}
	counter int
}

// NewInMemoryEmailRepository creates a new in-memory email repository
func NewInMemoryEmailRepository() *InMemoryEmailRepository {
	return &InMemoryEmailRepository{
		data: make(map[string]map[string]interface{}),
	}
}

// SaveEmailLog saves an email log entry
func (r *InMemoryEmailRepository) SaveEmailLog(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	id := fmt.Sprintf("email_%d", r.counter)
	r.data[id] = map[string]interface{}{
		"id":      id,
		"to":      to,
		"subject": subject,
		"body":    body,
		"status":  "sent",
	}
	return nil
}

// GetEmailLog retrieves an email log entry
func (r *InMemoryEmailRepository) GetEmailLog(ctx context.Context, id string) (map[string]interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if email, ok := r.data[id]; ok {
		return email, nil
	}
	return nil, fmt.Errorf("email log not found: %s", id)
}

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mu      sync.RWMutex
	data    map[string]map[string]interface{}
	counter int
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		data: make(map[string]map[string]interface{}),
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user map[string]interface{}) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	id := fmt.Sprintf("user_%d", r.counter)
	user["id"] = id
	r.data[id] = user
	return id, nil
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (map[string]interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if user, ok := r.data[id]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found: %s", id)
}

// Update updates a user
func (r *InMemoryUserRepository) Update(ctx context.Context, id string, user map[string]interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("user not found: %s", id)
	}

	user["id"] = id
	r.data[id] = user
	return nil
}

// Delete deletes a user
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return fmt.Errorf("user not found: %s", id)
	}

	delete(r.data, id)
	return nil
}

// InMemoryCounterRepository is an in-memory implementation of CounterRepository
type InMemoryCounterRepository struct {
	mu    sync.RWMutex
	value int
}

// NewInMemoryCounterRepository creates a new in-memory counter repository
func NewInMemoryCounterRepository() *InMemoryCounterRepository {
	return &InMemoryCounterRepository{
		value: 0,
	}
}

// GetCounter returns the current counter value
func (r *InMemoryCounterRepository) GetCounter(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.value, nil
}

// IncrementCounter increments the counter and returns the new value
func (r *InMemoryCounterRepository) IncrementCounter(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.value++
	return r.value, nil
}
