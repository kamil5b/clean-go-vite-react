package interfaces

import (
	"context"
)

type CounterRepository interface {
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)
}
