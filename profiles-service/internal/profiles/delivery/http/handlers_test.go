package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/middleware"
	"cyansnbrst/profiles-service/internal/models"
	mock_profiles "cyansnbrst/profiles-service/internal/profiles/mock"
)

func TestProfilesHandlers_GetInfo(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesUC := mock_profiles.NewMockUseCase(ctrl)
	mockKafkaWriter := &kafka.Writer{}
	profilesHandlers := NewProfilesHandlers(cfg, mockProfilesUC, logger, mockKafkaWriter)

	tests := []struct {
		name         string
		userUID      string
		mockBehavior func(mockProfilesUC *mock_profiles.MockUseCase)
		expectStatus int
	}{
		{
			name:    "success",
			userUID: "53345",
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {
				mockProfilesUC.EXPECT().Get("53345").Return(&models.Profile{UserUID: "53345"}, nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:    "not found",
			userUID: "532",
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {
				mockProfilesUC.EXPECT().Get("532").Return(nil, errors.New("not found"))
			},
			expectStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesUC)

			req := httptest.NewRequest(http.MethodGet, "/profiles", nil)
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, tt.userUID))

			rr := httptest.NewRecorder()
			profilesHandlers.GetInfo().ServeHTTP(rr, req)
			require.Equal(t, tt.expectStatus, rr.Code)
		})
	}
}

func TestProfilesHandlers_EditData(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesUC := mock_profiles.NewMockUseCase(ctrl)
	mockKafkaWriter := &kafka.Writer{}
	profilesHandlers := NewProfilesHandlers(cfg, mockProfilesUC, logger, mockKafkaWriter)

	newLocation := "Obninsk"

	tests := []struct {
		name         string
		userUID      string
		requestBody  string
		mockBehavior func(mockProfilesUC *mock_profiles.MockUseCase)
		expectStatus int
	}{
		{
			name:        "success",
			userUID:     "234",
			requestBody: `{"location":"Obninsk","interests":["music","sports"]}`,
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {
				mockProfilesUC.EXPECT().Update("234", &newLocation, []string{"music", "sports"}).Return(nil)
				mockProfilesUC.EXPECT().SendToKafka(gomock.Any(), gomock.Any(), gomock.Any(), mockKafkaWriter).Return(nil)
			},
			expectStatus: http.StatusOK,
		},
		{
			name:         "invalid JSON",
			userUID:      "54",
			requestBody:  `{"location":}`,
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesUC)

			req := httptest.NewRequest(http.MethodPut, "/profiles/edit", strings.NewReader(tt.requestBody))
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, tt.userUID))

			rr := httptest.NewRecorder()
			profilesHandlers.EditData().ServeHTTP(rr, req)
			require.Equal(t, tt.expectStatus, rr.Code)
		})
	}
}

func TestProfilesHandlers_CreateProfile(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesUC := mock_profiles.NewMockUseCase(ctrl)
	mockKafkaWriter := &kafka.Writer{}
	profilesHandlers := NewProfilesHandlers(cfg, mockProfilesUC, logger, mockKafkaWriter)

	tests := []struct {
		name         string
		userUID      string
		requestBody  string
		mockBehavior func(mockProfilesUC *mock_profiles.MockUseCase)
		expectStatus int
	}{
		{
			name:        "success",
			userUID:     "543",
			requestBody: `{"name":"user"}`,
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {
				mockProfilesUC.EXPECT().CreateProfile("543", "user").Return(nil)
			},
			expectStatus: http.StatusCreated,
		},
		{
			name:         "missing name",
			userUID:      "12345",
			requestBody:  `{}`,
			mockBehavior: func(mockProfilesUC *mock_profiles.MockUseCase) {},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesUC)

			req := httptest.NewRequest(http.MethodPost, "/profiles/create/"+tt.userUID, strings.NewReader(tt.requestBody))

			params := httprouter.Params{
				httprouter.Param{
					Key:   "uid",
					Value: tt.userUID,
				},
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)
			req = req.WithContext(context.WithValue(req.Context(), middleware.UserContextKey, tt.userUID))

			rr := httptest.NewRecorder()
			profilesHandlers.CreateProfile().ServeHTTP(rr, req)
			require.Equal(t, tt.expectStatus, rr.Code)
		})
	}

}
