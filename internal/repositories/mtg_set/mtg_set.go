package mtg_set

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MtgSetRepository struct {
	q *repo.Queries
}

func NewMTGSetRepository(q *repo.Queries) *MtgSetRepository {
	return &MtgSetRepository{q: q}
}

func (r *MtgSetRepository) CreateMTGSet(ctx context.Context, set models.MTGSet) (uuid.UUID, error) {
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
