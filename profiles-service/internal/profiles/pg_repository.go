package profiles

import "cyansnbrst/profiles-service/internal/models"

// Profiles repository interface
type Repository interface {
	Get(uid string) (*models.Profile, error)
	Update(*models.Profile) error
	CreateProfile(uid string, name string, defaultLocation string, defaultInterests []string) error
}
