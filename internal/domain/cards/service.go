package cards

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/apperrors"
	"github.com/google/uuid"
)

type CardsService interface {
	ListCards(ctx context.Context) ([]repo.Card, error)
	FindCardById(ctx context.Context, id uuid.UUID) (repo.Card, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) CardsService {
	return &svc{repo: repo}
}

func (s *svc) ListCards(ctx context.Context) ([]repo.Card, error) {
	return s.repo.ListCards(ctx)
}

func (s *svc) FindCardById(ctx context.Context, id uuid.UUID) (repo.Card, error) {
	card, err := s.repo.FindCardById(ctx, id)
	return card, apperrors.PgxErrors(err)
}
