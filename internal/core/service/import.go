package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type importService struct {
	source     port.ScryfallCardSource
	cards      port.CardRepository
	games      port.GameRepository
	mtgSetRepo port.MTGSetRepository
	importRepo port.ScryfallImportRepository
}

func NewImportService(
	source port.ScryfallCardSource,
	cards port.CardRepository,
	games port.GameRepository,
	mtgSetRepo port.MTGSetRepository,
	importRepo port.ScryfallImportRepository,
) port.ImportService {
	return &importService{
		source:     source,
		cards:      cards,
		games:      games,
		mtgSetRepo: mtgSetRepo,
		importRepo: importRepo,
	}
}

func (s *importService) Run(ctx context.Context, gameCode string, limit int) error {
	slog.Info("Finding game by code")
	game, err := s.games.FindGameByCode(ctx, gameCode)
	if err != nil {
		return err
	}

	slog.Info("Getting scryfall bulk data")
	bulkData, err := s.source.GetBulkData(ctx)
	if err != nil {
		return err
	}

	var allCards *domain.ScryfallBulkData
	for _, bd := range bulkData {
		if bd.Type == "all_cards" {
			allCards = &bd
		}
	}

	if allCards == nil {
		return fmt.Errorf("bulk data request failed")
	}

	scryfallImp, err := s.importRepo.GetScryfallImportCount(ctx, allCards.UpdatedAt)

	if scryfallImp > 0 {
		slog.Info("Bulk data already imported.")
		return nil
	}

	slog.Info("Creating scryfall import")
	newImportId, err := s.importRepo.CreateImport(ctx, domain.NewScryfallImport{
		StartedAt:     time.Now(),
		BulkUpdatedAt: allCards.UpdatedAt,
		Status:        domain.ImportStarted,
	})

	slog.Info("Getting existent bulk file")
	filePath, exists := s.source.GetBulkFileIfExists(allCards.UpdatedAt)

	if !exists {
		slog.Info("Downloading bulk file")
		filePath, err = s.source.Download(ctx, *allCards)
		if err != nil {
			return err
		}
	}

	slog.Info("Requesting and reading sets")
	sets, err := s.source.GetSets(ctx)
	if err != nil {
		return err
	}

	slog.Info("Reading cards")
	cards, err := s.source.ReadCards(ctx, filePath, limit)
	if err != nil {
		return err
	}

	slog.Info("Creating sets")
	setsIds := make(map[uuid.UUID]uuid.UUID, len(sets))

	for _, set := range sets {
		setId, err := s.mtgSetRepo.CreateMTGSet(ctx, set)
		if err != nil {
			return err
		}

		setsIds[set.ImportID] = setId
	}

	for i := range cards {
		cards[i].Card.GameID = game.ID
		cards[i].MTGCard.SetID = setsIds[cards[i].MTGCard.SetID]
		cards[i].MTGCard.LastImportId = newImportId
	}

	slog.Info("Creating cards")
	if err := s.cards.CreateCards(ctx, cards); err != nil {
		return err
	}

	slog.Info("Finishing import")
	err = s.importRepo.UpdateScryfallImport(ctx, domain.UpdateScryfallImport{
		ID:         newImportId,
		FinishedAt: time.Now(),
		Status:     domain.ImportFinished,
	})
	if err != nil {
		return err
	}

	slog.Info("import finished", "game", gameCode, "cards", len(cards))

	return nil
}
