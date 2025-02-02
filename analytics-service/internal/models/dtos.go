package models

// Kafka message DTO
type KafkaMessageDTO struct {
	Action string   `json:"action"`
	Time   string   `json:"time"`
	Tags   []string `json:"tags"`
}
