package domain

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type CardInCollection struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int16     `json:"quantity"`
}

type NewCollection struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type CardToCollection struct {
	CardID       uuid.UUID `json:"card_id"`
	CollectionID uuid.UUID `json:"collection_id"`
	Quantity     int16     `json:"quantity"`
}
