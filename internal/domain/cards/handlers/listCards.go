package cards_handler

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/json"
)

func (h *CardsHandler) ListCardsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListCards(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}
