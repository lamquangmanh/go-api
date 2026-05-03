// Package db provides database connection utilities for PostgreSQL using pgx/v5.
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a new pgxpool connection pool from a DSN string.
// Returns an error if the DSN is invalid or the database is unreachable.
// It pings the database after connecting to verify the connection is healthy.
func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	// Ping to ensure the database is reachable at startup
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
