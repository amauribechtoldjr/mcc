package apperror

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrBadRequest            = errors.New("bad request")
	ErrInternalServerError   = errors.New("internal server error")
	ErrResourceAlreadyExists = errors.New("resource already exists")
)
