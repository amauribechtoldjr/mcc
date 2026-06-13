package port

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/google/uuid"
)

type CardService interface {
	ListCards(ctx context.Context) ([]domain.Card, error)
	FindCardByID(ctx context.Context, id uuid.UUID) (domain.Card, error)
}

type CardRepository interface {
	ListCards(ctx context.Context) ([]domain.Card, error)
	FindCardByID(ctx context.Context, id uuid.UUID) (domain.Card, error)
}
