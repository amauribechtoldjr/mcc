package web

import (
	"encoding/json"
	"net/http"
)

type ErrorResult struct {
	Error string `json:"error"`
}

type EntityResult struct {
	Data any `json:"data"`
}

func WriteError(w http.ResponseWriter, err error) {
	setDefaultHeaders(w, HTTPStatus(err))

	json.NewEncoder(w).Encode(&ErrorResult{Error: err.Error()})
}

func Write(w http.ResponseWriter, status int, data any) {
	setDefaultHeaders(w, status)

	json.NewEncoder(w).Encode(&EntityResult{Data: data})
}

func Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func setDefaultHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
