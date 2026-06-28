package scryfallimport_usecases

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/models"
	"github.com/amauribechtoldjr/mcc/internal/repositories"
	"github.com/amauribechtoldjr/mcc/internal/sources"
	"github.com/google/uuid"
)

type ScryfallImportUseCases struct {
	source *sources.Sources
	repos  *repositories.Repositories
}

func NewImportUseCases(
	sources *sources.Sources,
	repos *repositories.Repositories,
) *ScryfallImportUseCases {
	return &ScryfallImportUseCases{
		source: sources,
		repos:  repos,
	}
}

func (s *ScryfallImportUseCases) Run(ctx context.Context, gameCode string, limit int) error {
	slog.Info("Finding game by code")
	game, err := s.repos.Game.FindGameByCode(ctx, gameCode)
	if err != nil {
		return err
	}

	slog.Info("Getting scryfall bulk data")
	bulkData, err := s.source.Scryfall.GetBulkData(ctx)
	if err != nil {
		return err
	}

	var allCards *models.ScryfallBulkData
	for _, bd := range bulkData {
		if bd.Type == "all_cards" {
			allCards = &bd
		}
	}

	if allCards == nil {
		return fmt.Errorf("bulk data request failed")
	}

	scryfallImp, err := s.repos.ScryfallImport.GetScryfallImportCount(ctx, allCards.UpdatedAt)

	if scryfallImp > 0 {
		slog.Info("Bulk data already imported.")
		return nil
	}

	slog.Info("Creating scryfall import")
	newImportId, err := s.repos.ScryfallImport.CreateImport(ctx, models.NewScryfallImport{
		StartedAt:     time.Now(),
		BulkUpdatedAt: allCards.UpdatedAt,
		Status:        models.ImportStarted,
	})

	slog.Info("Getting existent bulk file")
	filePath, exists := s.source.Scryfall.GetBulkFileIfExists(allCards.UpdatedAt)

	if !exists {
		slog.Info("Downloading bulk file")
		filePath, err = s.source.Scryfall.Download(ctx, *allCards)
		if err != nil {
			return fmt.Errorf("failed to download all cards bulk data: %w", err)
		}
	}

	slog.Info("Requesting and reading sets")
	sets, err := s.source.Scryfall.GetSets(ctx)
	if err != nil {
		return err
	}

	slog.Info("Reading cards")
	cards, err := s.source.Scryfall.ReadCards(ctx, filePath, limit)
	if err != nil {
		return err
	}

	slog.Info("Creating sets")
	setsIds := make(map[uuid.UUID]uuid.UUID, len(sets))

	for _, set := range sets {
		setId, err := s.repos.MtgSet.CreateMTGSet(ctx, set)
		if err != nil {
			return err
		}

		setsIds[set.ImportID] = setId
	}

	imported := make([]models.ImportCard, 0, len(cards))
	for i := range cards {
		setID, ok := setsIds[cards[i].MTGCard.SetID]
		if !ok {
			slog.Warn("skipping card: set not imported",
				"set_id", cards[i].MTGCard.SetID,
				"name", cards[i].MTGCard.Name)
			continue
		}

		cards[i].Card.GameID = game.ID
		cards[i].MTGCard.SetID = setID
		cards[i].MTGCard.LastImportId = newImportId
		imported = append(imported, cards[i])
	}

	slog.Info("Creating cards")
	if err := s.repos.Card.CreateCards(ctx, imported); err != nil {
		return err
	}

	slog.Info("Finishing import")
	err = s.repos.ScryfallImport.UpdateScryfallImport(ctx, models.UpdateScryfallImport{
		ID:         newImportId,
		FinishedAt: time.Now(),
		Status:     models.ImportFinished,
	})
	if err != nil {
		return err
	}

	slog.Info("import finished", "game", gameCode, "cards", len(imported))

	return nil
}
