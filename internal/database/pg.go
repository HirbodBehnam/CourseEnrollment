package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgresDatabase will connect to a postgres database
func NewPostgresDatabase(connectionString string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), connectionString)
}
