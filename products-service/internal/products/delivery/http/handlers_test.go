package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/models"
	mock_products "cyansnbrst/products-service/internal/products/mock"
	"cyansnbrst/products-service/pkg/db"
)

func TestProductsHandlers_Get(t *testing.T) {
	cfg := &config.Config{}

	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsUC := mock_products.NewMockUseCase(ctrl)
	kafkaWriter := &kafka.Writer{}

	productHandler := NewProductsHandlers(cfg, mockProductsUC, logger, kafkaWriter, kafkaWriter)

	tests := []struct {
		name         string
		id           string
		mockBehavior func(mockProductsUC *mock_products.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful get product",
			id:   "1",
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Get(int64(1)).Return(&models.Product{ID: 1, Name: "product", Tags: []string{"all"}}, nil)
				mockProductsUC.EXPECT().SendToKafka(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "product not found",
			id:   "2",
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Get(int64(2)).Return(nil, db.ErrRecordNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsUC)

			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.id, nil)

			params := httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: tt.id,
				},
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			productHandler.Get().ServeHTTP(rr, req)
			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestProductsHandlers_Create(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsUC := mock_products.NewMockUseCase(ctrl)
	kafkaWriter := &kafka.Writer{}

	productHandler := NewProductsHandlers(cfg, mockProductsUC, logger, kafkaWriter, kafkaWriter)

	tests := []struct {
		name         string
		requestBody  models.CreateProductDTO
		mockBehavior func(mockProductsUC *mock_products.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful create product",
			requestBody: models.CreateProductDTO{
				Name: "product",
				Tags: []string{"new"},
			},
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Create("product", []string{"new"}).Return(int64(1), nil)
				mockProductsUC.EXPECT().SendToKafka(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "empty name",
			requestBody: models.CreateProductDTO{
				Name: "",
				Tags: []string{"new"},
			},
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {},
			wantStatus:   http.StatusBadRequest,
		},
		{
			name: "nil tags",
			requestBody: models.CreateProductDTO{
				Name: "product",
				Tags: nil,
			},
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {},
			wantStatus:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsUC)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/products/create", bytes.NewReader(jsonBody))
			rr := httptest.NewRecorder()

			productHandler.Create().ServeHTTP(rr, req)
			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestProductsHandlers_Update(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsUC := mock_products.NewMockUseCase(ctrl)
	kafkaWriter := &kafka.Writer{}

	productHandler := NewProductsHandlers(cfg, mockProductsUC, logger, kafkaWriter, kafkaWriter)

	updatedName := "new name"

	tests := []struct {
		name         string
		id           string
		requestBody  models.UpdateProductDTO
		mockBehavior func(mockProductsUC *mock_products.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful update product",
			id:   "1",
			requestBody: models.UpdateProductDTO{
				Name: &updatedName,
				Tags: []string{"updated"},
			},
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Update(int64(1), &updatedName, []string{"updated"}).Return(nil)
				mockProductsUC.EXPECT().SendToKafka(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "product not found",
			id:   "2",
			requestBody: models.UpdateProductDTO{
				Name: &updatedName,
				Tags: []string{"updated"},
			},
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Update(int64(2), &updatedName, []string{"updated"}).Return(db.ErrRecordNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsUC)

			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/products/update/"+tt.id, bytes.NewReader(jsonBody))

			params := httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: tt.id,
				},
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			productHandler.Update().ServeHTTP(rr, req)
			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestProductsHandlers_Delete(t *testing.T) {
	cfg := &config.Config{}
	logger := zap.NewNop()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductsUC := mock_products.NewMockUseCase(ctrl)
	kafkaWriter := &kafka.Writer{}

	productHandler := NewProductsHandlers(cfg, mockProductsUC, logger, kafkaWriter, kafkaWriter)

	tests := []struct {
		name         string
		id           string
		mockBehavior func(mockProductsUC *mock_products.MockUseCase)
		wantStatus   int
	}{
		{
			name: "successful delete product",
			id:   "1",
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Delete(int64(1)).Return(nil)
				mockProductsUC.EXPECT().SendToKafka(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "product not found",
			id:   "2",
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {
				mockProductsUC.EXPECT().Delete(int64(2)).Return(db.ErrRecordNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:         "invalid id parameter",
			id:           "invalid",
			mockBehavior: func(mockProductsUC *mock_products.MockUseCase) {},
			wantStatus:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockProductsUC)

			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.id, nil)

			params := httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: tt.id,
				},
			}

			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			productHandler.Delete().ServeHTTP(rr, req)
			require.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
