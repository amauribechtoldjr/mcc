package domain

import (
	"github.com/google/uuid"
)

type Card struct {
	ID       uuid.UUID `json:"id"`
	OracleID uuid.UUID `json:"oracle_id"`
	GameID   uuid.UUID `json:"game_id"`
}

type MTGCard struct {
	ID             uuid.UUID `json:"id"`
	SetID          uuid.UUID `json:"set_id"`
	CardID         uuid.UUID `json:"card_id"`
	Name           string    `json:"name"`
	Layout         string    `json:"layout"`
	CMC            float32   `json:"cmc"`
	ColorIdentity  []string  `json:"color_identity"`
	ColorIndicator []string  `json:"color_indicator"`
	Colors         []string  `json:"colors"`
	ImgSmallURI    string    `json:"img_smal_uri"`
	ImgNormalURI   string    `json:"img_normal_uri"`
}

type ImportCard struct {
	Card    Card
	MTGCard MTGCard
}
