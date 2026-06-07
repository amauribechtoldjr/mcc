package cards

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/apperrors"
	"github.com/amauribechtoldjr/mcc/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CardsHandler struct {
	service CardsService
}

func NewHandler(service CardsService) *CardsHandler {
	return &CardsHandler{
		service: service,
	}
}

func (h *CardsHandler) ListCardsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListCards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *CardsHandler) FindCardById(w http.ResponseWriter, r *http.Request) {
	cardId, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		json.WriteError(w, apperrors.ErrBadRequest)
		return
	}

	products, err := h.service.FindCardById(r.Context(), cardId)
	if err != nil {
		json.WriteError(w, err)
		return
	}

	json.Write(w, http.StatusOK, products)
}
