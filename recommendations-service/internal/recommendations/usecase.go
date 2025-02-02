package recommendations

import "cyansnbrst/recommendations-service/internal/models"

type UseCase interface {
	GenerateRecommendationsForUser(userUID string, newInterests []string) error
	UpdateRecommendationsForProduct(productID int64, newTags []string) error
	GetRecommendationsForUser(userUID string) ([]models.Recommendation, error)
	IncrementPopularity(productID int64) error
	InsertProduct(productID int64, tags []string) error
	DeleteProduct(productID int64) error
}
