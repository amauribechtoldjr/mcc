package cards

import (
	"net/http"
	"strconv"

	"github.com/amauribechtoldjr/mcc/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListCardsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListCards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindCardById(w http.ResponseWriter, r *http.Request) {
	cardId, err := strconv.ParseInt(chi.URLParam(r, "cardId"), 0, 0)
	if err != nil {
		json.WriteError(w, err)
		return
	}

	products, err := h.service.FindCardById(r.Context(), cardId)
	if err != nil {
		json.WriteError(w, err)
		return
	}

	json.Write(w, http.StatusOK, products)
}
