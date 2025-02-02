package auth

import "cyansnbrst/auth-service/internal/models"

// Auth repository interface
type Repository interface {
	Insert(user *models.User) (string, error)
	GetByEmail(email string) (*models.User, error)
}
