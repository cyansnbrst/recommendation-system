package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"cyansnbrst/auth-service/config"
	mock_auth "cyansnbrst/auth-service/internal/auth/mock"
	"cyansnbrst/auth-service/internal/models"
)

func TestAuthUseCase_Create(t *testing.T) {
	cfg := &config.Config{
		Timeout: config.Timeout{
			Token: time.Hour,
		},
		SecretKey: "secret",
	}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, logger)

	tests := []struct {
		name         string
		email        string
		password     string
		mockBehavior func(mockAuthRepo *mock_auth.MockRepository)
		wantErr      bool
	}{
		{
			name:     "success",
			email:    "test@test.com",
			password: "test",
			mockBehavior: func(mockAuthRepo *mock_auth.MockRepository) {
				mockAuthRepo.EXPECT().Insert(gomock.Any()).Return("7543", nil)
			},
			wantErr: false,
		},
		{
			name:     "insert error",
			email:    "test@example.com",
			password: "test",
			mockBehavior: func(mockAuthRepo *mock_auth.MockRepository) {
				mockAuthRepo.EXPECT().Insert(gomock.Any()).Return("", errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAuthRepo)

			token, userID, err := authUC.Create(tt.email, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, token)
				require.Empty(t, userID)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)
				require.NotEmpty(t, userID)
			}
		})
	}
}

func TestAuthUseCase_GenerateAndValidateJWT(t *testing.T) {
	cfg := &config.Config{
		Timeout: config.Timeout{
			Token: time.Hour,
		},
		SecretKey: "secret",
	}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, logger)

	user := models.User{
		ID:      "5748",
		IsAdmin: true,
	}

	tests := []struct {
		name      string
		token     string
		setup     func() string
		wantUID   string
		wantAdmin bool
		wantErr   bool
	}{
		{
			name: "valid token",
			setup: func() string {
				token, _ := authUC.GenerateJWT(user)
				return token
			},
			wantUID:   "5748",
			wantAdmin: true,
			wantErr:   false,
		},
		{
			name:      "invalid token format",
			token:     "wrong token",
			setup:     func() string { return "" },
			wantUID:   "",
			wantAdmin: false,
			wantErr:   true,
		},
		{
			name: "wrong secret key",
			setup: func() string {
				authUCWrong := NewAuthUseCase(&config.Config{SecretKey: "wrongkey"}, nil, logger)
				token, _ := authUCWrong.GenerateJWT(user)
				return token
			},
			wantUID:   "",
			wantAdmin: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if tt.setup != nil {
				token = tt.setup()
			}

			uid, isAdmin, err := authUC.ValidateToken(token)

			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, uid)
				require.False(t, isAdmin)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantUID, uid)
				require.Equal(t, tt.wantAdmin, isAdmin)
			}
		})
	}
}

func TestAuthUseCase_ValidateCredentials(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthRepo := mock_auth.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, logger)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:           "12345",
		Email:        "test@test.com",
		PasswordHash: string(hashedPassword),
	}

	tests := []struct {
		name         string
		email        string
		password     string
		mockBehavior func(mockAuthRepo *mock_auth.MockRepository)
		wantErr      bool
	}{
		{
			name:     "valid credentials",
			email:    "test@test.com",
			password: "correctpassword",
			mockBehavior: func(mockAuthRepo *mock_auth.MockRepository) {
				mockAuthRepo.EXPECT().GetByEmail("test@test.com").Return(mockUser, nil)
			},
			wantErr: false,
		},
		{
			name:     "user not found",
			email:    "noname@test.com",
			password: "correctpassword",
			mockBehavior: func(mockAuthRepo *mock_auth.MockRepository) {
				mockAuthRepo.EXPECT().GetByEmail("noname@test.com").Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
		{
			name:     "invalid password",
			email:    "test@test.com",
			password: "test",
			mockBehavior: func(mockAuthRepo *mock_auth.MockRepository) {
				mockAuthRepo.EXPECT().GetByEmail("test@test.com").Return(mockUser, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockAuthRepo)

			user, err := authUC.ValidateCredentials(tt.email, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, user)
			} else {
				require.NoError(t, err)
				require.NotNil(t, user)
				require.Equal(t, mockUser.ID, user.ID)
			}
		})
	}
}
