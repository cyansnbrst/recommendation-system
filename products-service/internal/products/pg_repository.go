package products

import "cyansnbrst/products-service/internal/models"

// Products repository interface
type Repository interface {
	Create(name string, tags []string) (int64, error)
	Update(product *models.Product) error
	Delete(productID int64) error
	GetByID(productID int64) (*models.Product, error)
}
