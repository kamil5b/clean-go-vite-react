package message

import (
	"testing"
)

func TestNewMessageService(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewMessageService()

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
