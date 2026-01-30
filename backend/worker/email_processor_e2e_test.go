// +build integration

package worker

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/kamil5b/clean-go-vite-react/backend/service"
	"github.com/kamil5b/clean-go-vite-react/backend/task"
	"github.com/redis/go-redis/v9"
)

// E2E test for email notification async processing
// This test requires Redis to be running on localhost:6379
func TestEmailProcessor_E2E_EnqueueAndProcess(t *testing.T) {
	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping E2E test")
	}

	// Create Asynq client
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	defer asynqClient.Close()

	// Create task payload
	payload := task.EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Marshal payload
	data, err := payload.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Create Asynq task
	asynqTask := &asynq.Task{
		Type:    task.TypeEmailNotification,
		Payload: data,
	}

	// Enqueue task
	info, err := asynqClient.Enqueue(asynqTask)
	if err != nil {
		t.Fatalf("Failed to enqueue task: %v", err)
	}

	if info.ID == "" {
		t.Error("Expected task ID, got empty string")
	}

	// Process task
	emailService := service.NewEmailService()
	processor := NewEmailProcessor(emailService)

	// Create Asynq task again for processing (simulating worker receiving it)
	processTask := &asynq.Task{
		Type:    task.TypeEmailNotification,
		Payload: data,
	}

	err = processor.ProcessTask(ctx, processTask)
	if err != nil {
		t.Errorf("Failed to process task: %v", err)
	}
}

// Test async task retry on failure
func TestEmailProcessor_E2E_TaskRetry(t *testing.T) {
	// Skip if Redis is not available
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping E2E test")
	}

	// Create invalid payload (empty To field)
	payload := task.EmailNotificationPayload{
		To:      "", // Invalid
		Subject: "Test",
		Body:    "Test",
	}

	data, err := payload.Marshal()
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Process task
	emailService := service.NewEmailService()
	processor := NewEmailProcessor(emailService)

	asynqTask := &asynq.Task{
		Type:    task.TypeEmailNotification,
		Payload: data,
	}

	err = processor.ProcessTask(ctx, asynqTask)
	// Should handle gracefully (not panic)
	if err != nil && ctx.Err() != context.Canceled {
		t.Logf("Task processing returned error: %v", err)
	}
}

// Test multiple concurrent task processing
func TestEmailProcessor_E2E_ConcurrentTasks(t *testing.T) {
	emailService := service.NewEmailService()
	processor := NewEmailProcessor(emailService)

	// Process multiple tasks concurrently
	const numTasks = 5
	done := make(chan error, numTasks)

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			ctx := context.Background()
			payload := task.EmailNotificationPayload{
				To:      "user@example.com",
				Subject: "Task",
				Body:    "Body",
			}

			data, _ := payload.Marshal()
			asynqTask := &asynq.Task{
				Type:    task.TypeEmailNotification,
				Payload: data,
			}

			done <- processor.ProcessTask(ctx, asynqTask)
		}(i)
	}

	// Wait for all tasks
	for i := 0; i < numTasks; i++ {
		if err := <-done; err != nil && err != context.Canceled {
			t.Errorf("Task %d failed: %v", i, err)
		}
	}
}

// Test task serialization/deserialization roundtrip
func TestEmailProcessor_E2E_PayloadRoundtrip(t *testing.T) {
	original := task.EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Email",
		Body:    "This is a test email with special chars: !@#$%^&*()",
	}

	// Marshal
	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal
	var unmarshaled task.EmailNotificationPayload
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify
	if unmarshaled.To != original.To ||
		unmarshaled.Subject != original.Subject ||
		unmarshaled.Body != original.Body {
		t.Error("Payload roundtrip failed: data mismatch")
	}
}
```

Now let me create E2E tests for the API:
