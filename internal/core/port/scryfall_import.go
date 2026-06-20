package port

import (
	"context"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/google/uuid"
)

type ScryfallImportRepository interface {
	CreateImport(ctx context.Context, importData domain.NewScryfallImport) (uuid.UUID, error)
	GetScryfallImportCount(ctx context.Context, updated_at time.Time) (int64, error)
	UpdateScryfallImport(ctx context.Context, updateImport domain.UpdateScryfallImport) error
}
