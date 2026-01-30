package counter

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
)

func TestNewCounterService(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*gomock.Controller) *mock.MockCounterRepository
		expectNil bool
	}{
		{
			name: "should create service with valid repository",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockCounterRepository {
				return mock.NewMockCounterRepository(ctrl)
			},
			expectNil: false,
		},
		{
			name: "should return non-nil service instance",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockCounterRepository {
				return mock.NewMockCounterRepository(ctrl)
			},
			expectNil: false,
		},
		{
			name: "should initialize service properly",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockCounterRepository {
				return mock.NewMockCounterRepository(ctrl)
			},
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := tt.setupRepo(ctrl)
			svc := NewCounterService(mockRepo)

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
