package postgres

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type collectionRepository struct {
	q repo.Querier
}

func NewCollectionRepository(q repo.Querier) port.CollectionRepository {
	return &collectionRepository{q: q}
}

func (r *collectionRepository) CreateCollection(ctx context.Context, in domain.NewCollection) (domain.Collection, error) {
	row, err := r.q.CreateCollection(ctx, repo.CreateCollectionParams{
		UserID: in.UserID,
		Name:   in.Name,
	})
	if err != nil {
		return domain.Collection{}, mapError(err)
	}

	return toDomainCollection(row), nil
}

func (r *collectionRepository) AddCardToCollection(ctx context.Context, in domain.CardToCollection) error {
	return mapError(r.q.AddCardToCollection(ctx, repo.AddCardToCollectionParams{
		CardID:       in.CardID,
		CollectionID: in.CollectionID,
		Quantity:     in.Quantity,
	}))
}

func (r *collectionRepository) ListCollections(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {
	rows, err := r.q.ListCollections(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}

	collections := make([]domain.Collection, 0, len(rows))
	for _, row := range rows {
		collections = append(collections, toDomainCollection(row))
	}

	return collections, nil
}

func toDomainCollection(row repo.Collection) domain.Collection {
	return domain.Collection{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UserID:    row.UserID,
	}
}
