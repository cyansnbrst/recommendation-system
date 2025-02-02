package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"cyansnbrst/analytics-service/config"
)

// InitKafkaReader initializes a Kafka consumer for the specified topic
func InitKafkaReader(cfg *config.Config, topicKey string, groupID string) (*kafka.Reader, error) {
	topic, exists := cfg.Kafka.Topics[topicKey]
	if !exists {
		return nil, fmt.Errorf("topic key '%s' not found in configuration", topicKey)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   topic,
		GroupID: cfg.Kafka.GroupID,
	})

	return reader, nil
}

// ConsumeMessages listens for messages from the Kafka topic
func ConsumeMessages(ctx context.Context, reader *kafka.Reader, handler func(kafka.Message) error) error {
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("error reading message: %w", err)
		}

		if err := handler(m); err != nil {
			return fmt.Errorf("error handling message: %w", err)
		}
	}
}
