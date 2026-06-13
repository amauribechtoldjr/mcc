package web

import (
	"errors"
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
)

func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, apperror.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, apperror.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, apperror.ErrResourceAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
