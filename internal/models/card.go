package models

import "github.com/google/uuid"

type Card struct {
	ID     uuid.UUID `json:"id"`
	GameID uuid.UUID `json:"game_id"`
}
