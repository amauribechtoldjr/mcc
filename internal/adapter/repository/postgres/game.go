package postgres

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
)

type gameRepository struct {
	q repo.Querier
}

func NewGameRepository(q repo.Querier) port.GameRepository {
	return &gameRepository{q: q}
}

func (r *gameRepository) FindGameByCode(ctx context.Context, code string) (domain.Game, error) {
	row, err := r.q.FindGameByCode(ctx, code)
	if err != nil {
		return domain.Game{}, mapError(err)
	}

	return toDomainGame(row), nil
}

func toDomainGame(row repo.Game) domain.Game {
	return domain.Game{
		ID:   row.ID,
		Name: row.Name,
		Code: row.Code,
	}
}
