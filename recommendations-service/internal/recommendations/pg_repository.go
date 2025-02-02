package recommendations

import "cyansnbrst/recommendations-service/internal/models"

// Recommendations repository interface
type Repository interface {
	CreateRecommendation(user_uid string, product_id int64) error
	GetRecommendationsByUser(user_uid string) ([]models.Recommendation, error)
	InsertUser(user_uid string, interests []string) error
	InsertProduct(product_id int64, tags []string) error
	IncrementPopularity(product_id int64) error
	UpdateProductTags(product_id int64, tags []string) error
	UpdateUserInterests(user_uid string, interests []string) error
	FindProductsByTags(tag string) ([]int64, error)
	DeleteRecommendationsForProduct(productID int64) error
	GetAllUsers() ([]string, error)
	GetUserInterests(userUID string) ([]string, error)
	DeleteRecommendationsForUser(userUID string) error
	DeleteProduct(productID int64) error
}
