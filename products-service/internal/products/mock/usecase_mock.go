// Code generated by MockGen. DO NOT EDIT.
// Source: internal/products/usecase.go

// Package mock_products is a generated GoMock package.
package mock_products

import (
	context "context"
	models "cyansnbrst/products-service/internal/models"
	kafka "cyansnbrst/products-service/pkg/kafka"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kafka0 "github.com/segmentio/kafka-go"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUseCase) Create(name string, tags []string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, tags)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUseCaseMockRecorder) Create(name, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), name, tags)
}

// Delete mocks base method.
func (m *MockUseCase) Delete(id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUseCaseMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUseCase)(nil).Delete), id)
}

// Get mocks base method.
func (m *MockUseCase) Get(id int64) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUseCaseMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUseCase)(nil).Get), id)
}

// SendToKafka mocks base method.
func (m *MockUseCase) SendToKafka(ctx context.Context, key string, message kafka.KafkaMessage, writer *kafka0.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendToKafka", ctx, key, message, writer)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendToKafka indicates an expected call of SendToKafka.
func (mr *MockUseCaseMockRecorder) SendToKafka(ctx, key, message, writer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendToKafka", reflect.TypeOf((*MockUseCase)(nil).SendToKafka), ctx, key, message, writer)
}

// Update mocks base method.
func (m *MockUseCase) Update(id int64, name *string, tags []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", id, name, tags)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUseCaseMockRecorder) Update(id, name, tags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUseCase)(nil).Update), id, name, tags)
}
