package auth

import (
	"net/http"

	"cyansnbrst/auth-service/internal/models"
)

// Auth usecase interface
type UseCase interface {
	Create(email, password string) (string, string, error)
	ValidateToken(tokenString string) (string, bool, error)
	GenerateJWT(user models.User) (string, error)
	ValidateCredentials(email, password string) (*models.User, error)
	CreateProfile(uid string, name string, cookie *http.Cookie) error
}
