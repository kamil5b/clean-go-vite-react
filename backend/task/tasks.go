package task

import (
	"encoding/json"
)

// Task type constants
const (
	TypeEmailNotification = "email:notification"
	TypeDataExport        = "data:export"
)

// EmailNotificationPayload contains the data for email notification tasks
type EmailNotificationPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Marshal converts the payload to JSON bytes
func (p *EmailNotificationPayload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalEmailNotificationPayload decodes JSON bytes to payload
func UnmarshalEmailNotificationPayload(data []byte) (*EmailNotificationPayload, error) {
	var payload EmailNotificationPayload
	err := json.Unmarshal(data, &payload)
	return &payload, err
}

// DataExportPayload contains the data for data export tasks
type DataExportPayload struct {
	UserID    string `json:"user_id"`
	Format    string `json:"format"` // "csv", "json", "xlsx"
	ExportURL string `json:"export_url"`
}

// Marshal converts the payload to JSON bytes
func (p *DataExportPayload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalDataExportPayload decodes JSON bytes to payload
func UnmarshalDataExportPayload(data []byte) (*DataExportPayload, error) {
	var payload DataExportPayload
	err := json.Unmarshal(data, &payload)
	return &payload, err
}
