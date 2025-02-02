package consumers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/analytics-service/config"
	mock_analytics "cyansnbrst/analytics-service/internal/analytics/mock"
)

func TestKafkaMessageHandlers_HandleMessage(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnalyticsUC := mock_analytics.NewMockUseCase(ctrl)

	analyticsHandlers := NewKafkaMessageHandlers(cfg, mockAnalyticsUC, logger)

	tests := []struct {
		name         string
		message      kafka.Message
		mockBehavior func(mockAnalyticsUC *mock_analytics.MockUseCase)
		wantErr      bool
	}{
		{
			name: "valid message",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"click","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockAnalyticsUC *mock_analytics.MockUseCase) {
				mockAnalyticsUC.EXPECT().Insert("click", "1234", gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid JSON format",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte("invalid_json"),
			},
			mockBehavior: func(mockAnalyticsUC *mock_analytics.MockUseCase) {},
			wantErr:      true,
		},
		{
			name: "time parsing error",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"click","tags":["tag1"],"time":"invalid_time"}`),
			},
			mockBehavior: func(mockAnalyticsUC *mock_analytics.MockUseCase) {},
			wantErr:      true,
		},
		{
			name: "insert error",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"click","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockAnalyticsUC *mock_analytics.MockUseCase) {
				mockAnalyticsUC.EXPECT().Insert("click", "1234", gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAnalyticsUC)

			err := analyticsHandlers.HandleMessage(tt.message)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
