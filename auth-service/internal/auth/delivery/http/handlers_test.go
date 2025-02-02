package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/auth-service/config"
	mock_auth "cyansnbrst/auth-service/internal/auth/mock"
	"cyansnbrst/auth-service/internal/models"
	"cyansnbrst/auth-service/pkg/db"
)

func TestAuthHandlers_Register(t *testing.T) {
	cfg := &config.Config{
		Timeout: config.Timeout{
			Token:  time.Hour,
			Cookie: time.Hour,
		},
		SecretKey: "secret",
		Env:       "development",
	}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock_auth.NewMockUseCase(ctrl)

	authHandler := NewAuthHandlers(cfg, mockAuthUC, logger)

	tests := []struct {
		name         string
		body         models.RegisterUserDTO
		mockBehavior func(mockAuthUC *mock_auth.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful registration",
			body: models.RegisterUserDTO{
				Email:    "test@test.com",
				Password: "test",
				Name:     "user",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				token := "token"
				mockAuthUC.EXPECT().Create("test@test.com", "test").Return(token, "12345", nil)
				mockAuthUC.EXPECT().CreateProfile("12345", "user", gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "duplicate email error",
			body: models.RegisterUserDTO{
				Email:    "test@test.com",
				Password: "test",
				Name:     "user",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().Create("test@test.com", "test").Return("", "", db.ErrDuplicateEmail)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "creation error",
			body: models.RegisterUserDTO{
				Email:    "test@test.com",
				Password: "test",
				Name:     "user",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().Create("test@test.com", "test").Return("", "", errors.New("error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "profile creation error",
			body: models.RegisterUserDTO{
				Email:    "test@test.com",
				Password: "test",
				Name:     "user",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				token := "token"
				mockAuthUC.EXPECT().Create("test@test.com", "test").Return(token, "12345", nil)
				mockAuthUC.EXPECT().CreateProfile("12345", "user", gomock.Any()).Return(errors.New("profile creation error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAuthUC)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			authHandler.Register().ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)

			if tt.wantStatus == http.StatusOK {
				cookies := rr.Result().Cookies()
				require.NotEmpty(t, cookies)
			}
		})
	}
}

func TestAuthHandlers_Login(t *testing.T) {
	cfg := &config.Config{
		Timeout: config.Timeout{
			Token:  time.Hour,
			Cookie: time.Hour,
		},
		SecretKey: "secret",
		Env:       "development",
	}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock_auth.NewMockUseCase(ctrl)
	authHandler := NewAuthHandlers(cfg, mockAuthUC, logger)

	tests := []struct {
		name         string
		body         models.LoginUserDTO
		mockBehavior func(mockAuthUC *mock_auth.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful login",
			body: models.LoginUserDTO{
				Email:    "test@test.com",
				Password: "correctpassword",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().ValidateCredentials("test@test.com", "correctpassword").Return(&models.User{ID: "12345"}, nil)
				mockAuthUC.EXPECT().GenerateJWT(gomock.Any()).Return("token", nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: models.LoginUserDTO{
				Email:    "test@test.com",
				Password: "wrongpassword",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().ValidateCredentials("test@test.com", "wrongpassword").Return(nil, errors.New("invalid credentials"))
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "error generating JWT",
			body: models.LoginUserDTO{
				Email:    "test@test.com",
				Password: "correctpassword",
			},
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().ValidateCredentials("test@test.com", "correctpassword").Return(&models.User{ID: "12345"}, nil)
				mockAuthUC.EXPECT().GenerateJWT(gomock.Any()).Return("", errors.New("jwt error"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAuthUC)

			body, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			authHandler.Login().ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestAuthHandlers_TokenValidation(t *testing.T) {
	cfg := &config.Config{
		Timeout: config.Timeout{
			Token: time.Hour,
		},
		SecretKey: "secret",
		Env:       "development",
	}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthUC := mock_auth.NewMockUseCase(ctrl)
	authHandler := NewAuthHandlers(cfg, mockAuthUC, logger)

	tests := []struct {
		name         string
		token        string
		mockBehavior func(mockAuthUC *mock_auth.MockUseCase)
		wantStatus   int
	}{
		{
			name:  "valid token",
			token: "token",
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().ValidateToken("token").Return("12345", true, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "invalid token",
			token: "invalid_token",
			mockBehavior: func(mockAuthUC *mock_auth.MockUseCase) {
				mockAuthUC.EXPECT().ValidateToken("invalid_token").Return("", false, errors.New("invalid token"))
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAuthUC)

			req := httptest.NewRequest(http.MethodGet, "/auth/authenticate", nil)
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: tt.token,
			})
			rr := httptest.NewRecorder()

			authHandler.TokenValidation().ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
