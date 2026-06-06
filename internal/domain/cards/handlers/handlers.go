package cards_handler

import "github.com/amauribechtoldjr/mcc/internal/domain/cards"

type CardsHandler struct {
	service cards.Service
}

func NewHandler(service cards.Service) *CardsHandler {
	return &CardsHandler{
		service: service,
	}
}
