package analytics

import "github.com/segmentio/kafka-go"

// Analytics handlers interface
type KafkaHandlers interface {
	HandleMessage(msg kafka.Message) error
}
