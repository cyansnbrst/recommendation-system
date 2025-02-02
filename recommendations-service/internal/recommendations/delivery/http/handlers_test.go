package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/middleware"
	"cyansnbrst/recommendations-service/internal/models"
	mock_recommendations "cyansnbrst/recommendations-service/internal/recommendations/mock"
)

func TestRecommendationsHandlers_GetInfo(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRecommensationsUC := mock_recommendations.NewMockUseCase(ctrl)
	recommendationsHandlers := NewRecommendationsHandlers(cfg, mockRecommensationsUC, logger)

	tests := []struct {
		name         string
		userUID      string
		mockBehavior func(mockRecommensationsUC *mock_recommendations.MockUseCase)
		expectStatus int
	}{
		{
			name:    "success",
			userUID: "53345",
			mockBehavior: func(mockRecommensationsUC *mock_recommendations.MockUseCase) {
				mockRecommensationsUC.EXPECT().GetRecommendationsForUser("53345").Return([]models.Recommendation{{ID: 1, UserUID: "53345", ProductID: 1}}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:    "not found",
			userUID: "532",
			mockBehavior: func(mockRecommensationsUC *mock_recommendations.MockUseCase) {
				mockRecommensationsUC.EXPECT().GetRecommendationsForUser("532").Return(nil, errors.New("db error"))
			},
			expectStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRecommensationsUC)

			req := httptest.NewRequest(http.MethodGet, "/recommendations", nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, tt.userUID))

			rr := httptest.NewRecorder()
			recommendationsHandlers.GetInfo().ServeHTTP(rr, req)
			require.Equal(t, tt.expectStatus, rr.Code)
		})
	}
}
