package products

import (
	"context"

	"github.com/segmentio/kafka-go"

	"cyansnbrst/products-service/internal/models"
	kf "cyansnbrst/products-service/pkg/kafka"
)

// Products usecase interface
type UseCase interface {
	Get(id int64) (*models.Product, error)
	Update(id int64, name *string, tags []string) error
	Create(name string, tags []string) (int64, error)
	Delete(id int64) error
	SendToKafka(ctx context.Context, key string, message kf.KafkaMessage, writer *kafka.Writer) error
}
