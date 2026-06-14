package domain

import (
	"github.com/google/uuid"
)

type Card struct {
	ID       uuid.UUID `json:"id"`
	OracleID uuid.UUID `json:"oracle_id"`
	GameID   uuid.UUID `json:"game_id"`
}
