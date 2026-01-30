package email

import (
	"testing"
)

func TestNewEmailService(t *testing.T) {
	tests := []struct {
		name      string
		expectNil bool
	}{
		{
			name:      "should create service successfully",
			expectNil: false,
		},
		{
			name:      "should return non-nil service instance",
			expectNil: false,
		},
		{
			name:      "should initialize service properly",
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService()

			if tt.expectNil {
				if svc != nil {
					t.Errorf("expected nil service, got non-nil")
				}
			} else {
				if svc == nil {
					t.Errorf("expected non-nil service, got nil")
				}
			}
		})
	}
}
