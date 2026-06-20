package postgres

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type mtgSetRepository struct {
	q repo.Querier
}

func NewMTGSetRepository(q repo.Querier) port.MTGSetRepository {
	return &mtgSetRepository{q: q}
}

func (r *mtgSetRepository) CreateMTGSet(ctx context.Context, set domain.MTGSet) (uuid.UUID, error) {
	setId, err := r.q.CreateMTGSet(ctx, repo.CreateMTGSetParams{
		ImportID: set.ImportID,
		Name:     set.Name,
		Code:     set.Code,
		ReleasedAt: pgtype.Timestamptz{
			Time:  set.ReleasedAt,
			Valid: true,
		},
		ParentSetCode: &set.ParentSetCode,
		CardCount:     &set.CardCount,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return setId, nil
}
