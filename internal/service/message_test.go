package service

import (
	"context"
	"testing"
)

func TestMessageService_GetMessage(t *testing.T) {
	svc := NewMessageService()

	msg, err := svc.GetMessage(context.Background())
	if err != nil {
		t.Fatalf("GetMessage failed: %v", err)
	}

	expected := "Hello, from the golang World!"
	if msg != expected {
		t.Errorf("expected %q, got %q", expected, msg)
	}
}

func TestMessageService_GetMessageWithCanceledContext(t *testing.T) {
	svc := NewMessageService()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	msg, err := svc.GetMessage(ctx)
	// Service ignores context and returns message anyway
	if err != nil {
		t.Fatalf("GetMessage failed: %v", err)
	}
	if msg != "Hello, from the golang World!" {
		t.Errorf("expected message, got %q", msg)
	}
}
