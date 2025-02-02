package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/models"
	mock_recommendations "cyansnbrst/recommendations-service/internal/recommendations/mock"
)

func TestRecommendationsUC_GenerateRecommendationsForUser(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		userUID      string
		newInterests []string
		mockBehavior func(mockRepo *mock_recommendations.MockRepository)
		wantErr      bool
	}{
		{
			name:         "success new user",
			userUID:      "user1",
			newInterests: []string{"tag1", "tag2"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().GetUserInterests("user1").Return(nil, nil)
				mockRepo.EXPECT().InsertUser("user1", []string{"tag1", "tag2"}).Return(nil)
				mockRepo.EXPECT().FindProductsByTags("tag1").Return([]int64{1, 2}, nil)
				mockRepo.EXPECT().FindProductsByTags("tag2").Return([]int64{3}, nil)
				mockRepo.EXPECT().CreateRecommendation("user1", int64(1)).Return(nil)
				mockRepo.EXPECT().CreateRecommendation("user1", int64(2)).Return(nil)
				mockRepo.EXPECT().CreateRecommendation("user1", int64(3)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "success existing user",
			userUID:      "user2",
			newInterests: []string{"tag3"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().GetUserInterests("user2").Return([]string{"tag1"}, nil)
				mockRepo.EXPECT().DeleteRecommendationsForUser("user2").Return(nil)
				mockRepo.EXPECT().FindProductsByTags("tag3").Return([]int64{4}, nil)
				mockRepo.EXPECT().CreateRecommendation("user2", int64(4)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:         "error get user interests",
			userUID:      "user3",
			newInterests: []string{"tag1"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().GetUserInterests("user3").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := recommendationsUC.GenerateRecommendationsForUser(tt.userUID, tt.newInterests)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecommendationsUC_GetRecommendationsForUser(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		userUID      string
		mockBehavior func(mockRepo *mock_recommendations.MockRepository, mockRedisRepo *mock_recommendations.MockRedisRepository)
		wantErr      bool
	}{
		{
			name:    "success from cache",
			userUID: "user1",
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository, mockRedisRepo *mock_recommendations.MockRedisRepository) {
				mockRedisRepo.EXPECT().GetRecommendations("user1").Return([]models.Recommendation{{ProductID: 1}}, nil)
			},
			wantErr: false,
		},
		{
			name:    "success from db",
			userUID: "user2",
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository, mockRedisRepo *mock_recommendations.MockRedisRepository) {
				mockRedisRepo.EXPECT().GetRecommendations("user2").Return(nil, nil)
				mockRepo.EXPECT().GetRecommendationsByUser("user2").Return([]models.Recommendation{{ProductID: 2}}, nil)
				mockRedisRepo.EXPECT().SetRecommendations("user2", []models.Recommendation{{ProductID: 2}}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "db error",
			userUID: "user3",
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository, mockRedisRepo *mock_recommendations.MockRedisRepository) {
				mockRedisRepo.EXPECT().GetRecommendations("user3").Return(nil, nil)
				mockRepo.EXPECT().GetRecommendationsByUser("user3").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo, mockRedisRepo)

			_, err := recommendationsUC.GetRecommendationsForUser(tt.userUID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecommendationsUC_UpdateRecommendationsForProduct(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		productID    int64
		newTags      []string
		mockBehavior func(mockRepo *mock_recommendations.MockRepository)
		wantErr      bool
	}{
		{
			name:      "success",
			productID: 1,
			newTags:   []string{"tag1"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().GetAllUsers().Return([]string{"user1"}, nil)
				mockRepo.EXPECT().GetUserInterests("user1").Return([]string{"tag1"}, nil)
				mockRepo.EXPECT().DeleteRecommendationsForProduct(int64(1)).Return(nil)
				mockRepo.EXPECT().CreateRecommendation("user1", int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "error get all users",
			productID: 2,
			newTags:   []string{"tag2"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := recommendationsUC.UpdateRecommendationsForProduct(tt.productID, tt.newTags)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecommendationsUC_IncrementPopularity(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		productID    int64
		mockBehavior func(mockRepo *mock_recommendations.MockRepository)
		wantErr      bool
	}{
		{
			name:      "success",
			productID: 1,
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().IncrementPopularity(int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "db error",
			productID: 2,
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().IncrementPopularity(int64(2)).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := recommendationsUC.IncrementPopularity(tt.productID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecommendationsUC_InsertProduct(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		productID    int64
		tags         []string
		mockBehavior func(mockRepo *mock_recommendations.MockRepository)
		wantErr      bool
	}{
		{
			name:      "success",
			productID: 1,
			tags:      []string{"tag1", "tag2"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().InsertProduct(int64(1), []string{"tag1", "tag2"}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "db error",
			productID: 2,
			tags:      []string{"tag3"},
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().InsertProduct(int64(2), []string{"tag3"}).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := recommendationsUC.InsertProduct(tt.productID, tt.tags)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecommendationsUC_DeleteProduct(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_recommendations.NewMockRepository(ctrl)
	mockRedisRepo := mock_recommendations.NewMockRedisRepository(ctrl)
	recommendationsUC := NewRecommendationsUseCase(cfg, mockRepo, mockRedisRepo, logger)

	tests := []struct {
		name         string
		productID    int64
		mockBehavior func(mockRepo *mock_recommendations.MockRepository)
		wantErr      bool
	}{
		{
			name:      "success",
			productID: 1,
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().DeleteProduct(int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "db error",
			productID: 2,
			mockBehavior: func(mockRepo *mock_recommendations.MockRepository) {
				mockRepo.EXPECT().DeleteProduct(int64(2)).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			err := recommendationsUC.DeleteProduct(tt.productID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
