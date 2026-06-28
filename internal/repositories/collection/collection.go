package collection

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/db"
	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/models"

	"github.com/google/uuid"
)

type CollectionRepository struct {
	q *repo.Queries
}

func NewCollectionRepository(q *repo.Queries) *CollectionRepository {
	return &CollectionRepository{q: q}
}

func (r *CollectionRepository) CreateCollection(ctx context.Context, in models.NewCollection) (models.Collection, error) {
	row, err := r.q.CreateCollection(ctx, repo.CreateCollectionParams{
		UserID: in.UserID,
		Name:   in.Name,
	})
	if err != nil {
		return models.Collection{}, db.MapError(err)
	}

	return toDomainCollection(row), nil
}

func (r *CollectionRepository) AddCardToCollection(ctx context.Context, in models.CardToCollection) error {
	return db.MapError(r.q.AddCardToCollection(ctx, repo.AddCardToCollectionParams{
		CardID:       in.CardID,
		CollectionID: in.CollectionID,
		Quantity:     in.Quantity,
	}))
}

func (r *CollectionRepository) ListCollections(ctx context.Context, userID uuid.UUID) ([]models.Collection, error) {
	rows, err := r.q.ListCollections(ctx, userID)
	if err != nil {
		return nil, db.MapError(err)
	}

	collections := make([]models.Collection, 0, len(rows))
	for _, row := range rows {
		collections = append(collections, toDomainCollection(row))
	}

	return collections, nil
}

func toDomainCollection(row repo.Collection) models.Collection {
	return models.Collection{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UserID:    row.UserID,
	}
}
