package initialize

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	DB *sql.DB
}

func InitDB(ctx context.Context, config *Config) (*DB, error) {
	connStr := "postgres://" + config.PGUser + ":" + config.PGPassword + "@" + config.PGHost + ":" + config.PGPort + "/" + config.PGDatabase + "?sslmode=disable"
	pool, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := pool.PingContext(timeoutCtx); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(pool, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &DB{
		DB: pool,
	}, nil
}
