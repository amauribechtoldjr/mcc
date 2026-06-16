package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Connect opens a single pgx connection using the given DSN.
func Connect(ctx context.Context, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
