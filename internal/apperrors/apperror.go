package apperrors

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound              = errors.New("not found")
	ErrBadRequest            = errors.New("bad request")
	ErrInternalServerError   = errors.New("internal server error")
	ErrResourceAlreadyExists = errors.New("resource already exists")
)

func PgxErrors(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrResourceAlreadyExists
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return ErrNotFound
	default:
		return err
	}
}

func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, ErrResourceAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
