package port

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/google/uuid"
)

type CollectionService interface {
	CreateCollection(ctx context.Context, in domain.NewCollection) (domain.Collection, error)
	AddCardToCollection(ctx context.Context, in domain.CardToCollection) error
	ListCollections(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error)
}

type CollectionRepository interface {
	CreateCollection(ctx context.Context, in domain.NewCollection) (domain.Collection, error)
	AddCardToCollection(ctx context.Context, in domain.CardToCollection) error
	ListCollections(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error)
}
