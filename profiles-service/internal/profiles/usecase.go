package profiles

import (
	"context"

	"github.com/segmentio/kafka-go"

	"cyansnbrst/profiles-service/internal/models"
	kf "cyansnbrst/profiles-service/pkg/kafka"
)

// Profiles usecase interface
type UseCase interface {
	Get(uid string) (*models.Profile, error)
	Update(uid string, location *string, interests []string) error
	CreateProfile(uid string, name string) error
	SendToKafka(ctx context.Context, key string, message kf.KafkaMessage, writer *kafka.Writer) error
}
