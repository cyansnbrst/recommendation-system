package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	mock_analytics "cyansnbrst/analytics-service/internal/analytics/mock"
)

func TestAnalyticsUseCase_Insert(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyticsRepo := mock_analytics.NewMockRepository(ctrl)
	analyticsUC := NewAnalyticsUseCase(cfg, mockAnalyticsRepo, logger)

	tests := []struct {
		name         string
		action       string
		objectID     string
		actionTime   time.Time
		mockBehavior func(mockAnalyticsRepo *mock_analytics.MockRepository)
		wantErr      bool
	}{
		{
			name:       "success",
			action:     "click",
			objectID:   "1234",
			actionTime: time.Now(),
			mockBehavior: func(mockAnalyticsRepo *mock_analytics.MockRepository) {
				mockAnalyticsRepo.EXPECT().Insert("click", "1234", gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:       "insert error",
			action:     "view",
			objectID:   "5678",
			actionTime: time.Now(),
			mockBehavior: func(mockAnalyticsRepo *mock_analytics.MockRepository) {
				mockAnalyticsRepo.EXPECT().Insert("view", "5678", gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAnalyticsRepo)

			err := analyticsUC.Insert(tt.action, tt.objectID, tt.actionTime)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
