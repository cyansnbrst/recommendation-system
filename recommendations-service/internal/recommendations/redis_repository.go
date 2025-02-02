package recommendations

import "cyansnbrst/recommendations-service/internal/models"

type RedisRepository interface {
	GetRecommendations(key string) ([]models.Recommendation, error)
	SetRecommendations(key string, recommendations []models.Recommendation) error
}
