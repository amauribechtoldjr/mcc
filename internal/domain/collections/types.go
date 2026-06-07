package collections

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
)

type CollectionsService interface {
	CreateCollection(
		ctx context.Context,
		collectionData repo.CreateCollectionParams,
	) (repo.Collection, error)

	AddCardToCollection(ctx context.Context, cardCollectionData repo.AddCardToCollectionParams) error
}
