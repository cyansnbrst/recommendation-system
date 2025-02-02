package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"

	"cyansnbrst/profiles-service/config"
)

// Kafka message struct
type KafkaMessage struct {
	Action string   `json:"action"`
	Time   string   `json:"time"`
	Tags   []string `json:"tags"`
}

// Init kafka producer with given topic
func InitKafkaWriter(cfg *config.Config, topicKey string) (*kafka.Writer, error) {
	topic, exists := cfg.Kafka.Topics[topicKey]
	if !exists {
		return nil, fmt.Errorf("topic key '%s' not found in configuration", topicKey)
	}

	writer := &kafka.Writer{
		Addr:        kafka.TCP(cfg.Kafka.Brokers...),
		Topic:       topic,
		Balancer:    &kafka.LeastBytes{},
		MaxAttempts: cfg.Kafka.MaxAttempts,
	}

	return writer, nil
}
