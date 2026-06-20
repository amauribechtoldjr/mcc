package port

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
)

type CardSource interface {
	Download(ctx context.Context) (filePath string, err error)
	ReadCards(ctx context.Context, filePath string, limit int) ([]domain.ImportCard, error)
	ReadSets(ctx context.Context) ([]domain.MTGSet, error)
	GetBulkFileIfExists() (string, bool)
}

type ImportService interface {
	Run(ctx context.Context, gameCode string, limit int) error
}
