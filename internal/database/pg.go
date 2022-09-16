package database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

// NewPostgresDatabase will connect to a postgres database
func NewPostgresDatabase(connectionString string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), connectionString)
}
