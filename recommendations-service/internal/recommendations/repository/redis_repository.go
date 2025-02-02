package repository

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/models"
	"cyansnbrst/recommendations-service/internal/recommendations"
)

// Recommendations redis repository
type recommendationsRedisRepo struct {
	cfg         *config.Config
	redisClient *redis.Client
}

// Recommendations repository constructor
func NewRecommendationsRedisRepository(cfg *config.Config, redisClient *redis.Client) recommendations.RedisRepository {
	return &recommendationsRedisRepo{cfg: cfg, redisClient: redisClient}
}

// Get recommendations for user
func (r *recommendationsRedisRepo) GetRecommendations(key string) ([]models.Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.RedisAction)
	defer cancel()

	recommendationsBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var recommendations []models.Recommendation
	if err = json.Unmarshal(recommendationsBytes, &recommendations); err != nil {
		return nil, err
	}

	return recommendations, nil
}

// Cache recommendations for user
func (r *recommendationsRedisRepo) SetRecommendations(key string, recommendations []models.Recommendation) error {
	recommendationsBytes, err := json.Marshal(recommendations)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.RedisAction)
	defer cancel()

	if err = r.redisClient.Set(ctx, key, recommendationsBytes, r.cfg.Timeout.RedisCache).Err(); err != nil {
		return err
	}

	return nil
}
