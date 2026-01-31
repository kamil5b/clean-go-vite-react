package message

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
)

func TestNewMessageService(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*gomock.Controller) *mock.MockMessageRepository
		expectNil bool
	}{
		{
			name: "should create service with valid repository",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockMessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			expectNil: false,
		},
		{
			name: "should return non-nil service instance",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockMessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			expectNil: false,
		},
		{
			name: "should initialize service properly",
			setupRepo: func(ctrl *gomock.Controller) *mock.MockMessageRepository {
				return mock.NewMockMessageRepository(ctrl)
			},
			expectNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := tt.setupRepo(ctrl)
			svc := NewMessageService(mockRepo)

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
