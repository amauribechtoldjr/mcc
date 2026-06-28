package usecases

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/amauribechtoldjr/mcc/internal/sources"
	card_usecases "github.com/amauribechtoldjr/mcc/internal/usecases/card"
	collection_usecases "github.com/amauribechtoldjr/mcc/internal/usecases/collection"
	scryfallimport_usecases "github.com/amauribechtoldjr/mcc/internal/usecases/scryfall_import"
	"github.com/google/uuid"
)

type UseCases struct {
	Card interface {
		ListCards(ctx context.Context) ([]models.Card, error)
		FindCardByID(ctx context.Context, id uuid.UUID) (models.Card, error)
	}
	Collection interface {
		CreateCollection(ctx context.Context, in models.NewCollection) (models.Collection, error)
		AddCardToCollection(ctx context.Context, in models.CardToCollection) error
		ListCollections(ctx context.Context, userID uuid.UUID) ([]models.Collection, error)
	}
	ScryfallImport interface {
		Run(ctx context.Context, gameCode string, limit int) error
	}
}

func New(r *repositories.Repositories, s *sources.Sources) *UseCases {
	return &UseCases{
		Card:           card_usecases.NewCardUseCases(r),
		Collection:     collection_usecases.NewCollectionUseCases(r),
		ScryfallImport: scryfallimport_usecases.NewImportUseCases(s, r),
	}
}
