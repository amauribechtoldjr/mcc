package domain

import (
	"github.com/google/uuid"
)

type Card struct {
	ID     uuid.UUID `json:"id"`
	GameID uuid.UUID `json:"game_id"`
}

type MTGCard struct {
	ID              uuid.UUID `json:"id"` // = card.id (shared primary key)
	OracleID        uuid.UUID `json:"oracle_id"`
	SetID           uuid.UUID `json:"set_id"`
	SetCode         string    `json:"set"`
	Lang            string    `json:"lang"`
	CollectorNumber string    `json:"collector_number"`
	Name            string    `json:"name"` // printed_name se existir, senão o nome inglês
	PrintedTypeLine string    `json:"printed_type_line"`
	PrintedText     string    `json:"printed_text"`
	FlavorText      string    `json:"flavor_text"`
	Layout          string    `json:"layout"`
	CMC             float32   `json:"cmc"`
	ColorIdentity   []string  `json:"color_identity"`
	ColorIndicator  []string  `json:"color_indicator"`
	Colors          []string  `json:"colors"`
	ImgSmallURI     string    `json:"img_small_uri"`
	ImgNormalURI    string    `json:"img_normal_uri"`
	LastImportId    uuid.UUID `json:"last_import_id"`
}

type ImportCard struct {
	Card    Card
	MTGCard MTGCard
}
