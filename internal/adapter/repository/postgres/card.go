package postgres

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type cardRepository struct {
	q repo.Querier
}

func NewCardRepository(q repo.Querier) port.CardRepository {
	return &cardRepository{q: q}
}

func (r *cardRepository) ListCards(ctx context.Context) ([]domain.Card, error) {
	rows, err := r.q.ListCards(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	cards := make([]domain.Card, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, toDomainCard(row))
	}

	return cards, nil
}

func (r *cardRepository) FindCardByID(ctx context.Context, id uuid.UUID) (domain.Card, error) {
	row, err := r.q.FindCardById(ctx, id)
	if err != nil {
		return domain.Card{}, mapError(err)
	}

	return toDomainCard(row), nil
}

func toDomainCard(row repo.Card) domain.Card {
	return domain.Card{
		ID:       row.ID,
		OracleID: row.OracleID,
		GameID:   row.GameID,
	}
}
