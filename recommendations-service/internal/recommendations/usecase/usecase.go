package usecase

import (
	"go.uber.org/zap"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/models"
	"cyansnbrst/recommendations-service/internal/recommendations"
)

// Recommendations UseCase struct
type recommendationsUC struct {
	cfg                 *config.Config
	recommendationsRepo recommendations.Repository
	redisRepo           recommendations.RedisRepository
	logger              *zap.Logger
}

// New recommendations constructor
func NewRecommendationsUseCase(cfg *config.Config, recommendationsRepo recommendations.Repository, redisRepo recommendations.RedisRepository, logger *zap.Logger) recommendations.UseCase {
	return &recommendationsUC{cfg: cfg, recommendationsRepo: recommendationsRepo, redisRepo: redisRepo, logger: logger}
}

// Generate recommendations for user
func (u *recommendationsUC) GenerateRecommendationsForUser(userUID string, newInterests []string) error {
	interests, err := u.recommendationsRepo.GetUserInterests(userUID)
	if err != nil {
		return err
	}

	if interests != nil {
		return u.refreshRecommendations(userUID, newInterests)
	}

	err = u.recommendationsRepo.InsertUser(userUID, newInterests)
	if err != nil {
		return err
	}

	return u.createRecommendations(userUID, newInterests)
}

// Refresh existing recommendations
func (u *recommendationsUC) refreshRecommendations(userUID string, newInterests []string) error {
	err := u.recommendationsRepo.DeleteRecommendationsForUser(userUID)
	if err != nil {
		return err
	}

	return u.createRecommendations(userUID, newInterests)
}

// Create recommendations for user
func (u *recommendationsUC) createRecommendations(userUID string, interests []string) error {
	var productIDs []int64
	for _, interest := range interests {
		products, err := u.recommendationsRepo.FindProductsByTags(interest)
		if err != nil {
			return err
		}
		productIDs = append(productIDs, products...)
	}

	for _, productID := range productIDs {
		err := u.recommendationsRepo.CreateRecommendation(userUID, productID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update product tags and recommendations
func (u *recommendationsUC) UpdateRecommendationsForProduct(productID int64, newTags []string) error {
	users, err := u.recommendationsRepo.GetAllUsers()
	if err != nil {
		return err
	}

	err = u.recommendationsRepo.DeleteRecommendationsForProduct(productID)
	if err != nil {
		return err
	}

	for _, userUID := range users {
		userInterests, err := u.recommendationsRepo.GetUserInterests(userUID)
		if err != nil {
			return err
		}

		if isUserInterestedInNewTags(userInterests, newTags) {
			err := u.recommendationsRepo.CreateRecommendation(userUID, productID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Show user's recommendations
func (u *recommendationsUC) GetRecommendationsForUser(userUID string) ([]models.Recommendation, error) {
	recommendations, err := u.redisRepo.GetRecommendations(userUID)
	if err != nil {
		u.logger.Info("redis repository", zap.Error(err))
	}
	if recommendations != nil {
		u.logger.Info("got recommendations from the cache")
		return recommendations, nil
	}

	r, err := u.recommendationsRepo.GetRecommendationsByUser(userUID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetRecommendations(userUID, r); err != nil {
		u.logger.Error("redis repository", zap.Error(err))
	}

	return r, nil
}

// Increment product's popularity
func (u *recommendationsUC) IncrementPopularity(productID int64) error {
	return u.recommendationsRepo.IncrementPopularity(productID)
}

// Insert a new product
func (u *recommendationsUC) InsertProduct(productID int64, tags []string) error {
	return u.recommendationsRepo.InsertProduct(productID, tags)
}

// Delete product
func (u *recommendationsUC) DeleteProduct(productID int64) error {
	return u.recommendationsRepo.DeleteProduct(productID)
}

// Check if the user is interested in any of the new tags
func isUserInterestedInNewTags(userInterests []string, newTags []string) bool {
	for _, newTag := range newTags {
		for _, interest := range userInterests {
			if interest == newTag {
				return true
			}
		}
	}
	return false
}
