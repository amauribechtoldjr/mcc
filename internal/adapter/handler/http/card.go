package http

import (
	"net/http"

	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/amauribechtoldjr/mcc/internal/platform/apperror"
	"github.com/amauribechtoldjr/mcc/internal/platform/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CardHandler struct {
	service port.CardService
}

func NewCardHandler(service port.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) ListCards(w http.ResponseWriter, r *http.Request) {
	cards, err := h.service.ListCards(r.Context())
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusOK, cards)
}

func (h *CardHandler) FindCardByID(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		web.WriteError(w, apperror.ErrBadRequest)
		return
	}

	card, err := h.service.FindCardByID(r.Context(), cardID)
	if err != nil {
		web.WriteError(w, err)
		return
	}

	web.Write(w, http.StatusOK, card)
}
