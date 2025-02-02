package repository

import (
	"context"
	"database/sql"
	"time"

	"cyansnbrst/analytics-service/config"
	"cyansnbrst/analytics-service/internal/analytics"
)

// Analytics repository
type analyticsRepo struct {
	cfg *config.Config
	db  *sql.DB
}

// Analytics repository constructor
func NewAnalyticsRepository(cfg *config.Config, db *sql.DB) analytics.Repository {
	return &analyticsRepo{cfg: cfg, db: db}
}

// Insert a new analytics log
func (r *analyticsRepo) Insert(action string, objectID string, actionTime time.Time) error {
	query := `
		INSERT INTO actions (action, object_id, time)
		VALUES ($1, $2, $3)`

	args := []interface{}{action, objectID, actionTime}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
