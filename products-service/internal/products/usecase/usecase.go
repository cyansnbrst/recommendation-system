package usecase

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/models"
	"cyansnbrst/products-service/internal/products"
	kf "cyansnbrst/products-service/pkg/kafka"
)

// Products usecase struct
type productsUC struct {
	cfg          *config.Config
	productsRepo products.Repository
	logger       *zap.Logger
}

// Products usecase constructor
func NewProductsUseCase(cfg *config.Config, productsRepo products.Repository, logger *zap.Logger) products.UseCase {
	return &productsUC{cfg: cfg, productsRepo: productsRepo, logger: logger}
}

// Get a product by ID
func (u *productsUC) Get(id int64) (*models.Product, error) {
	return u.productsRepo.GetByID(id)
}

// Update a product
func (u *productsUC) Update(id int64, name *string, tags []string) error {
	product, err := u.productsRepo.GetByID(id)
	if err != nil {
		return err
	}

	if name != nil {
		product.Name = *name
	}
	if tags != nil {
		product.Tags = tags
	}

	return u.productsRepo.Update(product)
}

// Create a product
func (u *productsUC) Create(name string, tags []string) (int64, error) {
	return u.productsRepo.Create(name, tags)
}

// Delete a product
func (u *productsUC) Delete(id int64) error {
	return u.productsRepo.Delete(id)
}

// Send Kafka message
func (u *productsUC) SendToKafka(ctx context.Context, key string, message kf.KafkaMessage, writer *kafka.Writer) error {
	messageValue, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMessage := kafka.Message{
		Key:   []byte(key),
		Value: messageValue,
	}

	if err = writer.WriteMessages(ctx, kafkaMessage); err != nil {
		return err
	}

	return nil
}
