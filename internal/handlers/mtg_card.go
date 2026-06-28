package handlers

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
	"github.com/amauribechtoldjr/mcc/internal/platform/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handlers) registerMtgCardsEndpoints(r *chi.Mux) {
	r.Get("/mtg_cards", h.ListCards)
	r.Get("/mtg_cards/{cardId}", h.FindCardByID)
}

func (h *Handlers) ListCards(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) FindCardByID(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	card, err := h.useCases.Card.FindCardByID(r.Context(), cardID)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusOK, card)
}
