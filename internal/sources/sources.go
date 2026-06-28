package sources

import (
	"context"
	"net/http"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/models"
	scryfall_source "github.com/amauribechtoldjr/mcc/internal/sources/scryfall"
)

type Sources struct {
	Scryfall interface {
		Download(ctx context.Context, bulkData models.ScryfallBulkData) (string, error)
		ReadCards(ctx context.Context, filePath string, limit int) ([]models.ImportCard, error)
		GetSets(ctx context.Context) ([]models.MTGSet, error)
		GetBulkFileIfExists(updatedAt time.Time) (string, bool)
		GetBulkData(ctx context.Context) ([]models.ScryfallBulkData, error)
	}
}

func NewSource(client *http.Client) *Sources {
	return &Sources{
		Scryfall: scryfall_source.NewScryfallSource(client),
	}
}
