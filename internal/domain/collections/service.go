package collections

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/apperrors"
	"github.com/google/uuid"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) CollectionsService {
	return &svc{repo: repo}
}

func (s *svc) CreateCollection(ctx context.Context, collectionData repo.CreateCollectionParams) (repo.Collection, error) {
	newCollection, err := s.repo.CreateCollection(ctx, collectionData)
	return newCollection, apperrors.PgxErrors(err)
}

func (s *svc) AddCardToCollection(ctx context.Context, cardCollectionData repo.AddCardToCollectionParams) error {
	return apperrors.PgxErrors(s.repo.AddCardToCollection(ctx, cardCollectionData))
}

func (s *svc) ListCollectionCards(ctx context.Context, collectionId uuid.UUID) ([]repo.ListCollectionCardsRow, error) {
	return s.repo.ListCollectionCards(ctx, collectionId)
}
