package collection_usecases

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/google/uuid"
)

type CollectionUseCases struct {
	repo *repositories.Repositories
}

func NewCollectionUseCases(r *repositories.Repositories) *CollectionUseCases {
	return &CollectionUseCases{repo: r}
}

func (s *CollectionUseCases) CreateCollection(ctx context.Context, in models.NewCollection) (models.Collection, error) {
	return s.repo.Collection.CreateCollection(ctx, in)
}

func (s *CollectionUseCases) AddCardToCollection(ctx context.Context, in models.CardToCollection) error {
	return s.repo.Collection.AddCardToCollection(ctx, in)
}

func (s *CollectionUseCases) ListCollections(ctx context.Context, userID uuid.UUID) ([]models.Collection, error) {
	return s.repo.Collection.ListCollections(ctx, userID)
}
