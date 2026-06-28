package scryfall

import (
	"context"
	"time"

	repo "github.com/amauribechtoldjr/mcc/internal/db/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ScryfallRepository struct {
	q *repo.Queries
}

func NewScryfallImportRepository(q *repo.Queries) *ScryfallRepository {
	return &ScryfallRepository{q: q}
}

func (s *ScryfallRepository) CreateImport(ctx context.Context, importData models.NewScryfallImport) (uuid.UUID, error) {
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

func (s *ScryfallRepository) GetScryfallImportCount(ctx context.Context, updated_at time.Time) (int64, error) {
	return s.q.GetScryfallImportCount(ctx, pgtype.Timestamptz{
		Time:  updated_at,
		Valid: true,
	})
}

func (s *ScryfallRepository) UpdateScryfallImport(ctx context.Context, updateImport models.UpdateScryfallImport) error {
	return s.q.UpdateScryfallImport(ctx, repo.UpdateScryfallImportParams{
		ID: updateImport.ID,
		FinishedAt: pgtype.Timestamptz{
			Time:  updateImport.FinishedAt,
			Valid: true,
		},
		Status: string(updateImport.Status),
	})
}
