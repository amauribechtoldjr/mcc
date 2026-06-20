package port

import (
	"context"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
)

type ScryfallCardSource interface {
	Download(ctx context.Context, bulkData domain.ScryfallBulkData) (string, error)
	ReadCards(ctx context.Context, filePath string, limit int) ([]domain.ImportCard, error)
	GetSets(ctx context.Context) ([]domain.MTGSet, error)
	GetBulkFileIfExists(updatedAt time.Time) (string, bool)
	GetBulkData(ctx context.Context) ([]domain.ScryfallBulkData, error)
}

type ImportService interface {
	Run(ctx context.Context, gameCode string, limit int) error
}
