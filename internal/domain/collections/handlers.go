package collections

import (
	"log"
	"net/http"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/apperrors"
	"github.com/amauribechtoldjr/mcc/internal/json"
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
		json.WriteError(w, nil)
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
		json.WriteError(w, nil)
		return
	}

	json.Write(w, http.StatusCreated, nil)
}
