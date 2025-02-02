package usecase

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/models"
	mock_products "cyansnbrst/products-service/internal/products/mock"
)

func TestProductsUseCase_Get(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mock_products.NewMockRepository(ctrl)
	productsUC := NewProductsUseCase(cfg, mockProductsRepo, logger)

	tests := []struct {
		name         string
		mockBehavior func(mockRepo *mock_products.MockRepository)
		wantErr      bool
	}{
		{
			name: "get product success",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().GetByID(int64(1)).Return(&models.Product{ID: 1, Name: "product"}, nil)
			},
			wantErr: false,
		},
		{
			name: "get product not found",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().GetByID(int64(1)).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsRepo)

			_, err := productsUC.Get(1)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProductsUseCase_Update(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mock_products.NewMockRepository(ctrl)
	productsUC := NewProductsUseCase(cfg, mockProductsRepo, logger)

	tests := []struct {
		name         string
		mockBehavior func(mockRepo *mock_products.MockRepository)
		wantErr      bool
	}{
		{
			name: "update product success",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().GetByID(int64(1)).Return(&models.Product{ID: 1, Name: "name"}, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "update product repository error",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().GetByID(int64(1)).Return(&models.Product{ID: 1, Name: "name"}, nil)
				mockRepo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))
			},
			wantErr: true,
		},
		{
			name: "update product not found",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().GetByID(int64(1)).Return(nil, errors.New("get error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsRepo)

			newName := "newname"
			err := productsUC.Update(1, &newName, []string{"tag1", "tag2"})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProductsUseCase_Create(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mock_products.NewMockRepository(ctrl)
	productsUC := NewProductsUseCase(cfg, mockProductsRepo, logger)

	tests := []struct {
		name         string
		mockBehavior func(mockRepo *mock_products.MockRepository)
		wantErr      bool
	}{
		{
			name: "create product success",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().Create("product", []string{"tag1", "tag2"}).Return(int64(1), nil)
			},
			wantErr: false,
		},
		{
			name: "create product repository error",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().Create("product", []string{"tag1", "tag2"}).Return(int64(0), errors.New("create error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsRepo)

			_, err := productsUC.Create("product", []string{"tag1", "tag2"})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestProductsUseCase_Delete(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsRepo := mock_products.NewMockRepository(ctrl)
	productsUC := NewProductsUseCase(cfg, mockProductsRepo, logger)

	tests := []struct {
		name         string
		mockBehavior func(mockRepo *mock_products.MockRepository)
		wantErr      bool
	}{
		{
			name: "delete product success",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().Delete(int64(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "delete product repository error",
			mockBehavior: func(mockRepo *mock_products.MockRepository) {
				mockRepo.EXPECT().Delete(int64(1)).Return(errors.New("delete error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsRepo)

			err := productsUC.Delete(1)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
