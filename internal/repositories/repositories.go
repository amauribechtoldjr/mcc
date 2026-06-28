package repositories

import (
	"context"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/repositories/card"
	"github.com/amauribechtoldjr/mcc/internal/repositories/collection"
	"github.com/amauribechtoldjr/mcc/internal/repositories/game"
	"github.com/amauribechtoldjr/mcc/internal/repositories/mtg_set"
	"github.com/amauribechtoldjr/mcc/internal/repositories/scryfall"
	"github.com/google/uuid"
)

type Repositories struct {
	Card interface {
		ListCards(ctx context.Context) ([]models.Card, error)
		FindCardByID(ctx context.Context, id uuid.UUID) (models.Card, error)
		CreateCards(ctx context.Context, cards []models.ImportCard) error
	}
	Collection interface {
		CreateCollection(ctx context.Context, in models.NewCollection) (models.Collection, error)
		AddCardToCollection(ctx context.Context, in models.CardToCollection) error
		ListCollections(ctx context.Context, userID uuid.UUID) ([]models.Collection, error)
	}
	Game interface {
		FindGameByCode(ctx context.Context, code string) (models.Game, error)
	}
	MtgSet interface {
		CreateMTGSet(ctx context.Context, set models.MTGSet) (uuid.UUID, error)
	}
	ScryfallImport interface {
		CreateImport(ctx context.Context, importData models.NewScryfallImport) (uuid.UUID, error)
		GetScryfallImportCount(ctx context.Context, updated_at time.Time) (int64, error)
		UpdateScryfallImport(ctx context.Context, updateImport models.UpdateScryfallImport) error
	}
}

func NewRepositories(q *repo.Queries) *Repositories {
	return &Repositories{
		Card:           card.NewCardRepository(q),
		Collection:     collection.NewCollectionRepository(q),
		Game:           game.NewGameRepository(q),
		MtgSet:         mtg_set.NewMTGSetRepository(q),
		ScryfallImport: scryfall.NewScryfallImportRepository(q),
	}
}
