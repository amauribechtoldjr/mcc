package json

import (
	"encoding/json"
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/apperrors"
)

type ErrorResult struct {
	Error string `json:"error"`
}

type EntityResult struct {
	Data any `json:"data"`
}

func WriteError(w http.ResponseWriter, err error) {
	setDefaultHeaders(w, apperrors.HTTPStatus(err))

	json.NewEncoder(w).Encode(&ErrorResult{Error: err.Error()})
}

func Write(w http.ResponseWriter, status int, data any) {
	setDefaultHeaders(w, status)

	json.NewEncoder(w).Encode(&EntityResult{Data: data})
}

func setDefaultHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
