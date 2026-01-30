package repository

import (
	"context"
	"testing"
)

func TestInMemoryMessageRepository_GetMessage(t *testing.T) {
	repo := NewInMemoryMessageRepository()
	repo.data["default"] = "Hello World"

	msg, err := repo.GetMessage(context.Background())
	if err != nil {
		t.Fatalf("GetMessage failed: %v", err)
	}

	if msg != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", msg)
	}
}

func TestInMemoryMessageRepository_GetMessage_NotFound(t *testing.T) {
	repo := NewInMemoryMessageRepository()

	_, err := repo.GetMessage(context.Background())
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestInMemoryMessageRepository_GetMessage_ContextCanceled(t *testing.T) {
	repo := NewInMemoryMessageRepository()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.GetMessage(ctx)
	if err == nil {
		t.Error("expected context error, got nil")
	}
}

func TestInMemoryEmailRepository_SaveEmailLog(t *testing.T) {
	repo := NewInMemoryEmailRepository()

	err := repo.SaveEmailLog(context.Background(), "test@example.com", "Subject", "Body")
	if err != nil {
		t.Fatalf("SaveEmailLog failed: %v", err)
	}

	if len(repo.data) != 1 {
		t.Errorf("expected 1 email, got %d", len(repo.data))
	}
}

func TestInMemoryEmailRepository_GetEmailLog(t *testing.T) {
	repo := NewInMemoryEmailRepository()
	repo.SaveEmailLog(context.Background(), "test@example.com", "Subject", "Body")

	email, err := repo.GetEmailLog(context.Background(), "email_1")
	if err != nil {
		t.Fatalf("GetEmailLog failed: %v", err)
	}

	if email["to"] != "test@example.com" {
		t.Errorf("expected 'test@example.com', got %v", email["to"])
	}
}

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewInMemoryUserRepository()

	id, err := repo.Create(context.Background(), map[string]interface{}{
		"email": "user@example.com",
		"name":  "Test User",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if id == "" {
		t.Error("expected non-empty id")
	}
}

func TestInMemoryUserRepository_FindByID(t *testing.T) {
	repo := NewInMemoryUserRepository()
	id, _ := repo.Create(context.Background(), map[string]interface{}{
		"email": "user@example.com",
	})

	user, err := repo.FindByID(context.Background(), id)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if user["email"] != "user@example.com" {
		t.Errorf("expected 'user@example.com', got %v", user["email"])
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := NewInMemoryUserRepository()
	id, _ := repo.Create(context.Background(), map[string]interface{}{
		"email": "user@example.com",
	})

	err := repo.Update(context.Background(), id, map[string]interface{}{
		"email": "newemail@example.com",
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	user, _ := repo.FindByID(context.Background(), id)
	if user["email"] != "newemail@example.com" {
		t.Errorf("expected 'newemail@example.com', got %v", user["email"])
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := NewInMemoryUserRepository()
	id, _ := repo.Create(context.Background(), map[string]interface{}{
		"email": "user@example.com",
	})

	err := repo.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.FindByID(context.Background(), id)
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestInMemoryUserRepository_Delete_NotFound(t *testing.T) {
	repo := NewInMemoryUserRepository()

	err := repo.Delete(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error, got nil")
	}
}
