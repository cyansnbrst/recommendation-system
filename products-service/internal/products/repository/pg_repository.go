package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"cyansnbrst/products-service/config"
	"cyansnbrst/products-service/internal/models"
	"cyansnbrst/products-service/internal/products"
	"cyansnbrst/products-service/pkg/db"
)

// Products repository
type productsRepo struct {
	cfg *config.Config
	db  *sql.DB
}

// Products repository constructor
func NewProductsRepository(cfg *config.Config, db *sql.DB) products.Repository {
	return &productsRepo{cfg: cfg, db: db}
}

// Insert a new product
func (r *productsRepo) Create(name string, tags []string) (int64, error) {
	query := `
		INSERT INTO products (name, tags)
		VALUES ($1, $2)
		RETURNING id`

	args := []interface{}{
		name,
		pq.Array(tags),
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	var id int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Edit an existing product
func (r *productsRepo) Update(product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, tags = $2, version = version + 1
		WHERE id = $3 AND version = $4`

	args := []interface{}{
		product.Name,
		pq.Array(product.Tags),
		product.ID,
		product.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return db.ErrRecordNotFound
	}

	return nil
}

// Delete an existing product
func (r *productsRepo) Delete(productID int64) error {
	query := `
        DELETE FROM products
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		var exists bool
		checkQuery := `
			SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`
		err = r.db.QueryRowContext(ctx, checkQuery, productID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return db.ErrRecordNotFound
		}
		return db.ErrEditConflict
	}

	return nil
}

// Get product by ID
func (r *productsRepo) GetByID(productID int64) (*models.Product, error) {
	queryGet := `
        SELECT id, name, tags, version
        FROM products
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	product := &models.Product{}
	args := []interface{}{
		&product.ID,
		&product.Name,
		pq.Array(&product.Tags),
		&product.Version,
	}

	row := r.db.QueryRowContext(ctx, queryGet, productID)
	if err := row.Scan(args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrRecordNotFound
		}
		return nil, err
	}

	return product, nil
}
