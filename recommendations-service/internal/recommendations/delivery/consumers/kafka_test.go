package consumers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	mock_recommendations "cyansnbrst/recommendations-service/internal/recommendations/mock"
)

func TestKafkaMessageHandlers_HandleProductMessage(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationsUC := mock_recommendations.NewMockUseCase(ctrl)
	kafkaHandlers := NewKafkaMessageHandlers(cfg, mockRecommendationsUC, logger)

	tests := []struct {
		name         string
		message      kafka.Message
		mockBehavior func(mockRecommendationsUC *mock_recommendations.MockUseCase)
		wantErr      bool
	}{
		{
			name: "valid product create message",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"product_create","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().InsertProduct(int64(1234), []string{"tag1"}).Return(nil)
				mockRecommendationsUC.EXPECT().UpdateRecommendationsForProduct(int64(1234), []string{"tag1"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid JSON format",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte("invalid_json"),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {},
			wantErr:      true,
		},
		{
			name: "product create error",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"product_create","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().InsertProduct(int64(1234), []string{"tag1"}).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "valid product update message",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"product_update","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().UpdateRecommendationsForProduct(int64(1234), []string{"tag1"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "valid product delete message",
			message: kafka.Message{
				Key:   []byte("1234"),
				Value: []byte(`{"action":"product_delete","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().DeleteProduct(int64(1234)).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRecommendationsUC)

			err := kafkaHandlers.HandleProductMessage(tt.message)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKafkaMessageHandlers_HandleUserMessage(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommendationsUC := mock_recommendations.NewMockUseCase(ctrl)
	kafkaHandlers := NewKafkaMessageHandlers(cfg, mockRecommendationsUC, logger)

	tests := []struct {
		name         string
		message      kafka.Message
		mockBehavior func(mockRecommendationsUC *mock_recommendations.MockUseCase)
		wantErr      bool
	}{
		{
			name: "valid user update message",
			message: kafka.Message{
				Key:   []byte("user1234"),
				Value: []byte(`{"action":"user_update","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().GenerateRecommendationsForUser("user1234", []string{"tag1"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid JSON format",
			message: kafka.Message{
				Key:   []byte("user1234"),
				Value: []byte("invalid_json"),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {},
			wantErr:      true,
		},
		{
			name: "user update error",
			message: kafka.Message{
				Key:   []byte("user1234"),
				Value: []byte(`{"action":"user_update","tags":["tag1"],"time":"2023-01-01T12:00:00Z"}`),
			},
			mockBehavior: func(mockRecommendationsUC *mock_recommendations.MockUseCase) {
				mockRecommendationsUC.EXPECT().GenerateRecommendationsForUser("user1234", []string{"tag1"}).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRecommendationsUC)

			err := kafkaHandlers.HandleUserMessage(tt.message)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
