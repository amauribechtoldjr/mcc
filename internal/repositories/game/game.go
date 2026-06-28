package game

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/db"
	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/models"
)

type GameRepository struct {
	q *repo.Queries
}

func NewGameRepository(q *repo.Queries) *GameRepository {
	return &GameRepository{q: q}
}

func (r *GameRepository) FindGameByCode(ctx context.Context, code string) (models.Game, error) {
	row, err := r.q.FindGameByCode(ctx, code)
	if err != nil {
		return models.Game{}, db.MapError(err)
	}

	return toDomainGame(row), nil
}

func toDomainGame(row repo.Game) models.Game {
	return models.Game{
		ID:   row.ID,
		Name: row.Name,
		Code: row.Code,
	}
}
