package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"cyansnbrst/profiles-service/config"
)

// Open new postgres connection
func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.SSLMode,
	)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.PostgreSQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.PostgreSQL.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.PostgreSQL.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout.PostgreSQLConn)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
