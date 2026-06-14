package http

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
	"github.com/amauribechtoldjr/mcc/internal/platform/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CollectionHandler struct {
	service port.CollectionService
}

func NewCollectionHandler(service port.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}

func (h *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var in domain.NewCollection
	if err := web.Read(r, &in); err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	collection, err := h.service.CreateCollection(r.Context(), in)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusCreated, collection)
}

func (h *CollectionHandler) AddCardToCollection(w http.ResponseWriter, r *http.Request) {
	var in domain.CardToCollection
	if err := web.Read(r, &in); err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	if err := h.service.AddCardToCollection(r.Context(), in); err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusCreated, nil)
}

func (h *CollectionHandler) ListCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	collections, err := h.service.ListCollections(r.Context(), userID)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusOK, collections)
}
