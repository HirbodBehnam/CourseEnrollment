package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// NewPostgresDatabase will connect to a postgres database
func NewPostgresDatabase(connectionString string) (*pgxpool.Pool, error) {
	// Create pool
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	// Test DB
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	err = pool.Ping(ctx)
	cancel()
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
