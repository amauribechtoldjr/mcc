package db

import (
	"errors"

	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return apperror.ErrResourceAlreadyExists
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return apperror.ErrNotFound
	default:
		return err
	}
}
