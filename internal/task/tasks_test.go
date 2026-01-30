package task

import (
	"encoding/json"
	"testing"
)

func TestEmailNotificationPayload_Marshal(t *testing.T) {
	// Arrange
	payload := &EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	// Act
	data, err := payload.Marshal()

	// Assert
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if data == nil {
		t.Error("expected non-nil data")
	}

	// Verify it's valid JSON
	var result EmailNotificationPayload
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.To != payload.To {
		t.Errorf("expected To=%q, got %q", payload.To, result.To)
	}
}

func TestEmailNotificationPayload_UnmarshalEmailNotificationPayload(t *testing.T) {
	// Arrange
	original := &EmailNotificationPayload{
		To:      "recipient@example.com",
		Subject: "Welcome",
		Body:    "Welcome to our service",
	}
	data, _ := original.Marshal()

	// Act
	payload, err := UnmarshalEmailNotificationPayload(data)

	// Assert
	if err != nil {
		t.Fatalf("UnmarshalEmailNotificationPayload failed: %v", err)
	}

	if payload.To != original.To {
		t.Errorf("expected To=%q, got %q", original.To, payload.To)
	}
	if payload.Subject != original.Subject {
		t.Errorf("expected Subject=%q, got %q", original.Subject, payload.Subject)
	}
	if payload.Body != original.Body {
		t.Errorf("expected Body=%q, got %q", original.Body, payload.Body)
	}
}

func TestEmailNotificationPayload_UnmarshalEmailNotificationPayload_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte("invalid json")

	// Act
	_, err := UnmarshalEmailNotificationPayload(invalidJSON)

	// Assert
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestEmailNotificationPayload_RoundTrip(t *testing.T) {
	// Arrange
	original := &EmailNotificationPayload{
		To:      "user@example.com",
		Subject: "Important",
		Body:    "This is important",
	}

	// Act
	marshaled, _ := original.Marshal()
	unmarshaled, _ := UnmarshalEmailNotificationPayload(marshaled)

	// Assert
	if unmarshaled.To != original.To ||
		unmarshaled.Subject != original.Subject ||
		unmarshaled.Body != original.Body {
		t.Error("round trip failed")
	}
}

func TestDataExportPayload_Marshal(t *testing.T) {
	// Arrange
	payload := &DataExportPayload{
		UserID:    "user123",
		Format:    "csv",
		ExportURL: "s3://bucket/export.csv",
	}

	// Act
	data, err := payload.Marshal()

	// Assert
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result DataExportPayload
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.UserID != payload.UserID {
		t.Errorf("expected UserID=%q, got %q", payload.UserID, result.UserID)
	}
	if result.Format != payload.Format {
		t.Errorf("expected Format=%q, got %q", payload.Format, result.Format)
	}
}

func TestDataExportPayload_UnmarshalDataExportPayload(t *testing.T) {
	// Arrange
	original := &DataExportPayload{
		UserID:    "user456",
		Format:    "json",
		ExportURL: "s3://bucket/export.json",
	}
	data, _ := original.Marshal()

	// Act
	payload, err := UnmarshalDataExportPayload(data)

	// Assert
	if err != nil {
		t.Fatalf("UnmarshalDataExportPayload failed: %v", err)
	}

	if payload.UserID != original.UserID {
		t.Errorf("expected UserID=%q, got %q", original.UserID, payload.UserID)
	}
	if payload.Format != original.Format {
		t.Errorf("expected Format=%q, got %q", original.Format, payload.Format)
	}
}

func TestDataExportPayload_UnmarshalDataExportPayload_InvalidJSON(t *testing.T) {
	// Arrange
	invalidJSON := []byte("{invalid json}")

	// Act
	_, err := UnmarshalDataExportPayload(invalidJSON)

	// Assert
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDataExportPayload_RoundTrip(t *testing.T) {
	// Arrange
	original := &DataExportPayload{
		UserID:    "user789",
		Format:    "xlsx",
		ExportURL: "s3://bucket/export.xlsx",
	}

	// Act
	marshaled, _ := original.Marshal()
	unmarshaled, _ := UnmarshalDataExportPayload(marshaled)

	// Assert
	if unmarshaled.UserID != original.UserID ||
		unmarshaled.Format != original.Format ||
		unmarshaled.ExportURL != original.ExportURL {
		t.Error("round trip failed")
	}
}

func TestTaskTypeConstants(t *testing.T) {
	// Assert
	if TypeEmailNotification != "email:notification" {
		t.Errorf("expected TypeEmailNotification to be 'email:notification', got %q", TypeEmailNotification)
	}

	if TypeDataExport != "data:export" {
		t.Errorf("expected TypeDataExport to be 'data:export', got %q", TypeDataExport)
	}
}

func TestEmailNotificationPayload_EmptyFields(t *testing.T) {
	// Arrange
	payload := &EmailNotificationPayload{
		To:      "",
		Subject: "",
		Body:    "",
	}

	// Act
	data, err := payload.Marshal()

	// Assert
	if err != nil {
		t.Fatalf("Marshal failed for empty fields: %v", err)
	}

	var result EmailNotificationPayload
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
}

func TestEmailNotificationPayload_SpecialCharacters(t *testing.T) {
	// Arrange
	payload := &EmailNotificationPayload{
		To:      "user+test@example.com",
		Subject: "Subject with 中文 and émojis 🎉",
		Body:    "Body with special chars: !@#$%^&*()",
	}

	// Act
	data, err := payload.Marshal()

	// Assert
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	unmarshaled, err := UnmarshalEmailNotificationPayload(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if unmarshaled.To != payload.To {
		t.Errorf("To mismatch")
	}
	if unmarshaled.Subject != payload.Subject {
		t.Errorf("Subject mismatch")
	}
	if unmarshaled.Body != payload.Body {
		t.Errorf("Body mismatch")
	}
}
