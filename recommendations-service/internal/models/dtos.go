package models

// Kafka message DTO
type KafkaMessageDTO struct {
	Action string   `json:"action"`
	Time   string   `json:"time"`
	Tags   []string `json:"tags"`
}

// Recommendations response
type RecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error"`
}
