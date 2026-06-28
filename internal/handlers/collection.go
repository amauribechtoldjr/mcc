package handlers

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
	"github.com/amauribechtoldjr/mcc/internal/platform/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) registerCollectionsEndpoints(r *chi.Mux) {
	r.Post("/collections", h.CreateCollection)
	r.Get("/collections/{userId}", h.ListCollections)
	r.Post("/collections/cards", h.AddCardToCollection)
}

func (h *Handlers) ListCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	collections, err := h.useCases.Collection.ListCollections(r.Context(), userID)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusOK, collections)
}

func (h *Handlers) AddCardToCollection(w http.ResponseWriter, r *http.Request) {
	var in models.CardToCollection
	if err := web.Read(r, &in); err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	if err := h.useCases.Collection.AddCardToCollection(r.Context(), in); err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusCreated, nil)
}

func (h *Handlers) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var in models.NewCollection
	if err := web.Read(r, &in); err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	collection, err := h.useCases.Collection.CreateCollection(r.Context(), in)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusCreated, collection)
}
