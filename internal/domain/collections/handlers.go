package collections

import (
	"log"
	"net/http"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/apperrors"
	"github.com/amauribechtoldjr/mcc/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CollectionsHandlers struct {
	service CollectionsService
}

func NewHandler(service CollectionsService) *CollectionsHandlers {
	return &CollectionsHandlers{
		service: service,
	}
}

func (h *CollectionsHandlers) CreateCollection(w http.ResponseWriter, r *http.Request) {
	var tempCollection repo.CreateCollectionParams

	if err := json.Read(r, &tempCollection); err != nil {
		log.Println(err)
		json.WriteError(w, apperrors.ErrBadRequest)
		return
	}

	collectionData, err := h.service.CreateCollection(r.Context(), tempCollection)
	if err != nil {
		log.Println(err)
		json.WriteError(w, err)
		return
	}

	json.Write(w, http.StatusCreated, collectionData)
}

func (h *CollectionsHandlers) AddCardToCollection(w http.ResponseWriter, r *http.Request) {
	var tempCardCollection repo.AddCardToCollectionParams

	if err := json.Read(r, &tempCardCollection); err != nil {
		log.Println(err)
		json.WriteError(w, apperrors.ErrBadRequest)
		return
	}

	err := h.service.AddCardToCollection(r.Context(), tempCardCollection)
	if err != nil {
		log.Println(err)
		json.WriteError(w, err)
		return
	}

	json.Write(w, http.StatusCreated, nil)
}

func (h *CollectionsHandlers) ListCollectionCards(w http.ResponseWriter, r *http.Request) {
	collectionId, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		json.WriteError(w, apperrors.ErrBadRequest)
		return
	}

	cardsList, err := h.service.ListCollectionCards(r.Context(), collectionId)
	if err != nil {
		log.Println(err)
		json.WriteError(w, apperrors.ErrInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, cardsList)
}
