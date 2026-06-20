package postgres

import (
	"context"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type scryfallImports struct {
	q repo.Querier
}

func NewScryfallImportRepository(q repo.Querier) port.ScryfallImportRepository {
	return &scryfallImports{q: q}
}

func (s *scryfallImports) CreateImport(ctx context.Context, importData domain.NewScryfallImport) (uuid.UUID, error) {
	return s.q.CreateScryfallImport(ctx, repo.CreateScryfallImportParams{
		StartedAt: pgtype.Timestamptz{
			Time:  importData.StartedAt,
			Valid: true,
		},
		BulkUpdatedAt: pgtype.Timestamptz{
			Time:  importData.BulkUpdatedAt,
			Valid: true,
		},
		Status: string(importData.Status),
	})
}

func (s *scryfallImports) GetScryfallImportCount(ctx context.Context, updated_at time.Time) (int64, error) {
	return s.q.GetScryfallImportCount(ctx, pgtype.Timestamptz{
		Time:  updated_at,
		Valid: true,
	})
}

func (s *scryfallImports) UpdateScryfallImport(ctx context.Context, updateImport domain.UpdateScryfallImport) error {
	return s.q.UpdateScryfallImport(ctx, repo.UpdateScryfallImportParams{
		ID: updateImport.ID,
		FinishedAt: pgtype.Timestamptz{
			Time:  updateImport.FinishedAt,
			Valid: true,
		},
		Status: string(updateImport.Status),
	})
}
