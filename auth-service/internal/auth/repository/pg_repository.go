package repository

import (
	"context"
	"database/sql"
	"errors"

	"cyansnbrst/auth-service/config"
	"cyansnbrst/auth-service/internal/auth"
	"cyansnbrst/auth-service/internal/models"
	"cyansnbrst/auth-service/pkg/db"
)

// Auth repository
type authRepo struct {
	cfg *config.Config
	db  *sql.DB
}

// Auth repository constructor
func NewAuthRepository(cfg *config.Config, db *sql.DB) auth.Repository {
	return &authRepo{cfg: cfg, db: db}
}

// Insert a new user
func (r *authRepo) Insert(user *models.User) (string, error) {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at`

	args := []interface{}{user.Email, user.PasswordHash}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return "", db.ErrDuplicateEmail
		default:
			return "", err
		}
	}

	return user.ID, nil
}

func (r *authRepo) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, created_at, email, is_admin, password_hash
		FROM users
		WHERE email = $1`

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Email,
		&user.IsAdmin,
		&user.PasswordHash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, db.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
