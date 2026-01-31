package csrf

import (
	"testing"
)

func TestNewCSRFService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create csrf service successfully",
		},
		{
			name: "should return non-nil service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewCSRFService()

			if service == nil {
				t.Errorf("expected non-nil service, got nil")
			}

			_, ok := service.(*csrfService)
			if !ok {
				t.Errorf("expected *csrfService, got %T", service)
			}
		})
	}
}
