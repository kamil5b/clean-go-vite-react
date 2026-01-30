package platform

import (
	"github.com/hibiken/asynq"
)

// NewAsynqClient creates a new Asynq client for enqueueing tasks
func NewAsynqClient(cfg *AsynqConfig) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.RedisAddr,
	})
}

// NewAsynqServer creates a new Asynq server for processing tasks
func NewAsynqServer(cfg *AsynqConfig) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: cfg.RedisAddr,
		},
		asynq.Config{
			Concurrency: cfg.Concurrency,
		},
	)
}

// NewAsynqMux creates a new task multiplexer for handling different task types
func NewAsynqMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}
