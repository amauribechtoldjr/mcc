package domain

import (
	"time"

	"github.com/google/uuid"
)

type MTGSet struct {
	ID            uuid.UUID `json:"id"`
	ImportID      uuid.UUID `json:"import_id"`
	Code          string    `json:"code"`
	Name          string    `json:"name"`
	ReleasedAt    time.Time `json:"released_at"`
	ParentSetCode string    `json:"parent_set_dunk"`
	CardCount     int32     `json:"card_count"`
}
