package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

	"cyansnbrst/profiles-service/config"
	"cyansnbrst/profiles-service/internal/models"
	"cyansnbrst/profiles-service/internal/profiles"
	"cyansnbrst/profiles-service/pkg/db"
)

// Profiles repository
type profilesRepo struct {
	cfg *config.Config
	db  *sql.DB
}

// Profiles repository constructor
func NewProfilesRepository(cfg *config.Config, db *sql.DB) profiles.Repository {
	return &profilesRepo{cfg: cfg, db: db}
}

// Get profile info by UID
func (r *profilesRepo) Get(uid string) (*models.Profile, error) {
	query := `
		SELECT user_uid, name, location, interests
		FROM profiles
		WHERE user_uid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	var profile models.Profile
	err := r.db.QueryRowContext(ctx, query, uid).Scan(
		&profile.UserUID,
		&profile.Name,
		&profile.Location,
		pq.Array(&profile.Interests),
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, db.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &profile, nil
}

// Update profile data (location and interests)
func (r *profilesRepo) Update(profile *models.Profile) error {
	query := `
		UPDATE profiles
		SET location = $1, interests = $2
		WHERE user_uid = $3`

	args := []interface{}{
		profile.Location,
		pq.Array(profile.Interests),
		profile.UserUID,
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

func (r *profilesRepo) CreateProfile(uid string, name string, defaultLocation string, defaultInterests []string) error {
	query := `
		INSERT INTO profiles (user_uid, name, location, interests)
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{
		uid,
		name,
		defaultLocation,
		pq.Array(defaultInterests),
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.cfg.Timeout.PostgreSQLAction)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
