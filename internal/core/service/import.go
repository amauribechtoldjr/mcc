package service

import (
	"context"
	"log/slog"

	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type importService struct {
	source     port.CardSource
	cards      port.CardRepository
	games      port.GameRepository
	mtgSetRepo port.MTGSetRepository
}

func NewImportService(
	source port.CardSource,
	cards port.CardRepository,
	games port.GameRepository,
	mtgSetRepo port.MTGSetRepository,
) port.ImportService {
	return &importService{
		source:     source,
		cards:      cards,
		games:      games,
		mtgSetRepo: mtgSetRepo,
	}
}

func (s *importService) Run(ctx context.Context, gameCode string, limit int) error {
	slog.Info("Finding game by code")
	game, err := s.games.FindGameByCode(ctx, gameCode)
	if err != nil {
		return err
	}

	slog.Info("Getting existent bulk file")
	filePath, exists := s.source.GetBulkFileIfExists()

	if !exists {
		slog.Info("Downloading bulk file")
		filePath, err = s.source.Download(ctx)
		if err != nil {
			return err
		}
	}

	slog.Info("Requesting and reading sets")
	sets, err := s.source.ReadSets(ctx)
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
	}

	slog.Info("Creating cards")
	if err := s.cards.CreateCards(ctx, cards); err != nil {
		return err
	}

	slog.Info("import finished", "game", gameCode, "cards", len(cards))

	return nil
}
