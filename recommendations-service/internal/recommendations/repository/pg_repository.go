package repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"cyansnbrst/recommendations-service/config"
	"cyansnbrst/recommendations-service/internal/models"
	"cyansnbrst/recommendations-service/internal/recommendations"
)

// Recommendations repository
type recommendationsRepo struct {
	cfg *config.Config
	db  *sql.DB
}

// Recommendations repository constructor
func NewRecommendationsRepository(cfg *config.Config, db *sql.DB) recommendations.Repository {
	return &recommendationsRepo{cfg: cfg, db: db}
}

// Insert a new recommendation
func (r *recommendationsRepo) CreateRecommendation(userUID string, productID int64) error {
	query := `
		INSERT INTO recommendations (user_uid, product_id)
		VALUES ($1, $2)`

	args := []interface{}{userUID, productID}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Get recommendations for user
func (r *recommendationsRepo) GetRecommendationsByUser(userUID string) ([]models.Recommendation, error) {
	query := `
        SELECT r.product_id
        FROM recommendations r
        JOIN products p ON r.product_id = p.product_id
        WHERE r.user_uid = $1
        ORDER BY p.popularity DESC`

	var recommendations []models.Recommendation

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recommendation models.Recommendation
		if err := rows.Scan(&recommendation.ProductID); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, recommendation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recommendations, nil
}

// Insert new user
func (r *recommendationsRepo) InsertUser(userUID string, interests []string) error {
	query := `
        INSERT INTO users (user_uid, interests)
        VALUES ($1, $2)`

	args := []interface{}{userUID, pq.Array(interests)}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Insert new product
func (r *recommendationsRepo) InsertProduct(productID int64, tags []string) error {
	query := `
        INSERT INTO products (product_id, tags)
        VALUES ($1, $2)`

	args := []interface{}{productID, pq.Array(tags)}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Increment popularity
func (r *recommendationsRepo) IncrementPopularity(productID int64) error {
	query := `
        UPDATE products
        SET popularity = popularity + 1
        WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	return nil
}

// Update product tags
func (r *recommendationsRepo) UpdateProductTags(productID int64, tags []string) error {
	query := `
        UPDATE products
        SET tags = $1
        WHERE product_id = $2`

	args := []interface{}{pq.Array(tags), productID}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update user interests
func (r *recommendationsRepo) UpdateUserInterests(userUID string, interests []string) error {
	query := `
        UPDATE users
        SET interests = $1
        WHERE user_uid = $2`

	args := []interface{}{pq.Array(interests), userUID}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// Find products by tags
func (r *recommendationsRepo) FindProductsByTags(tag string) ([]int64, error) {
	query := `
        SELECT product_id
        FROM products
        WHERE $1 = ANY(tags)`

	var productIDs []int64

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int64
		if err := rows.Scan(&productID); err != nil {
			return nil, err
		}
		productIDs = append(productIDs, productID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productIDs, nil
}

// Delete recommendations for a product
func (r *recommendationsRepo) DeleteRecommendationsForProduct(productID int64) error {
	query := `
        DELETE FROM recommendations
        WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	return nil
}

// Get all users
func (r *recommendationsRepo) GetAllUsers() ([]string, error) {
	query := `
        SELECT user_uid
        FROM users`

	var users []string

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userUID string
		if err := rows.Scan(&userUID); err != nil {
			return nil, err
		}
		users = append(users, userUID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Get user interests
func (r *recommendationsRepo) GetUserInterests(userUID string) ([]string, error) {
	query := `
        SELECT interests
        FROM users
        WHERE user_uid = $1`

	var interests []string

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, userUID)
	if err := row.Scan(pq.Array(&interests)); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return interests, nil
}

// Delete recommendations for a user
func (r *recommendationsRepo) DeleteRecommendationsForUser(userUID string) error {
	query := `
        DELETE FROM recommendations
        WHERE user_uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, userUID)
	if err != nil {
		return err
	}

	return nil
}

// Delete product
func (r *recommendationsRepo) DeleteProduct(productID int64) error {
	query := `
        DELETE FROM products
        WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return err
	}

	return nil
}
