package card_usecases

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/google/uuid"
)

type CardUseCases struct {
	repo *repositories.Repositories
}

func NewCardUseCases(r *repositories.Repositories) *CardUseCases {
	return &CardUseCases{repo: r}
}

func (s *CardUseCases) ListCards(ctx context.Context) ([]models.Card, error) {
	return s.repo.Card.ListCards(ctx)
}

func (s *CardUseCases) FindCardByID(ctx context.Context, id uuid.UUID) (models.Card, error) {
	return s.repo.Card.FindCardByID(ctx, id)
}
