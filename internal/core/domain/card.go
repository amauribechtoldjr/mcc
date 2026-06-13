package domain

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SrcURL    string    `json:"src_url"`
	CreatedAt time.Time `json:"created_at"`
}
