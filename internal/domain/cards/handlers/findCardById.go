package cards_handler

import (
	"net/http"
	"strconv"

	"github.com/amauribechtoldjr/mcc/internal/json"
	"github.com/go-chi/chi/v5"
)

func (h *CardsHandler) FindCardById(w http.ResponseWriter, r *http.Request) {
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
