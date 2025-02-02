package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/models"
	mock_profiles "cyansnbrst/profiles-service/internal/profiles/mock"
)

func TestProfilesUseCase_Get(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesRepo := mock_profiles.NewMockRepository(ctrl)
	profilesUC := NewProfilesUseCase(cfg, mockProfilesRepo, logger)

	tests := []struct {
		name         string
		uid          string
		mockBehavior func(mockProfilesRepo *mock_profiles.MockRepository)
		wantErr      bool
	}{
		{
			name: "success",
			uid:  "53453",
			mockBehavior: func(mockProfilesRepo *mock_profiles.MockRepository) {
				mockProfilesRepo.EXPECT().Get("53453").Return(&models.Profile{UserUID: "53453"}, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			uid:  "654",
			mockBehavior: func(mockProfilesRepo *mock_profiles.MockRepository) {
				mockProfilesRepo.EXPECT().Get("654").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesRepo)
			profile, err := profilesUC.Get(tt.uid)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, profile)
			} else {
				require.NoError(t, err)
				require.NotNil(t, profile)
				require.Equal(t, tt.uid, profile.UserUID)
			}
		})
	}
}

func TestProfilesUseCase_Update(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesRepo := mock_profiles.NewMockRepository(ctrl)
	profilesUC := NewProfilesUseCase(cfg, mockProfilesRepo, logger)

	newLocaton := "Obninsk"

	tests := []struct {
		name         string
		uid          string
		location     *string
		interests    []string
		mockBehavior func(mockProfilesRepo *mock_profiles.MockRepository)
		wantErr      bool
	}{
		{
			name:      "success",
			uid:       "12345",
			location:  &newLocaton,
			interests: []string{"music", "sports"},
			mockBehavior: func(mockProfilesRepo *mock_profiles.MockRepository) {
				profile := &models.Profile{UserUID: "12345"}
				mockProfilesRepo.EXPECT().Get("12345").Return(profile, nil)
				mockProfilesRepo.EXPECT().Update(&models.Profile{UserUID: "12345", Location: newLocaton, Interests: []string{"music", "sports"}}).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "profile not found",
			uid:       "67890",
			location:  &newLocaton,
			interests: []string{"reading"},
			mockBehavior: func(mockProfilesRepo *mock_profiles.MockRepository) {
				mockProfilesRepo.EXPECT().Get("67890").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesRepo)
			err := profilesUC.Update(tt.uid, tt.location, tt.interests)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProfilesUseCase_CreateProfile(t *testing.T) {
	cfg := &config.Config{
		DefaultLocation:  "Moscow",
		DefaultInterests: "all",
	}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfilesRepo := mock_profiles.NewMockRepository(ctrl)
	profilesUC := NewProfilesUseCase(cfg, mockProfilesRepo, logger)

	tests := []struct {
		name         string
		uid          string
		uname        string
		mockBehavior func(mockProfilesRepo *mock_profiles.MockRepository)
		wantErr      bool
	}{
		{
			name:  "success",
			uid:   "12345",
			uname: "user",
			mockBehavior: func(mockProfilesRepo *mock_profiles.MockRepository) {
				mockProfilesRepo.EXPECT().CreateProfile("12345", "user", "Moscow", []string{"all"}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProfilesRepo)
			err := profilesUC.CreateProfile(tt.uid, tt.uname)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
