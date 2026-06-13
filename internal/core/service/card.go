package service

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type cardService struct {
	repo port.CardRepository
}

func NewCardService(repo port.CardRepository) port.CardService {
	return &cardService{repo: repo}
}

func (s *cardService) ListCards(ctx context.Context) ([]domain.Card, error) {
	return s.repo.ListCards(ctx)
}

func (s *cardService) FindCardByID(ctx context.Context, id uuid.UUID) (domain.Card, error) {
	return s.repo.FindCardByID(ctx, id)
}
