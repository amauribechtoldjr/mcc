package models

import (
	"time"

	"github.com/google/uuid"
)

type ScryfallImportStatus string

const (
	ImportStarted  ScryfallImportStatus = "started"
	ImportFailed   ScryfallImportStatus = "failed"
	ImportFinished ScryfallImportStatus = "finished"
)

type NewScryfallImport struct {
	StartedAt     time.Time
	BulkUpdatedAt time.Time
	Status        ScryfallImportStatus
}

type UpdateScryfallImport struct {
	ID         uuid.UUID
	FinishedAt time.Time
	Status     ScryfallImportStatus
}

type ScryfallBulkData struct {
	Type        string    `json:"type"`
	UpdatedAt   time.Time `json:"updated_at"`
	DonwloadURI string    `json:"download_uri"`
}

type ImportCard struct {
	Card    Card
	MTGCard NewMTGCard
}
