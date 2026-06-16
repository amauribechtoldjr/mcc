package service

import (
	"context"
	"log/slog"

	"github.com/amauribechtoldjr/mcc/internal/core/port"
)

type importService struct {
	source port.CardSource
	cards  port.CardRepository
	games  port.GameRepository
}

func NewImportService(
	source port.CardSource,
	cards port.CardRepository,
	games port.GameRepository,
) port.ImportService {
	return &importService{source: source, cards: cards, games: games}
}

func (s *importService) Run(ctx context.Context, gameCode string, limit int) error {
	game, err := s.games.FindGameByCode(ctx, gameCode)
	if err != nil {
		return err
	}

	filePath, exists := s.source.GetBulkFileIfExists()

	if !exists {
		filePath, err = s.source.Download(ctx)
		if err != nil {
			return err
		}
	}

	cards, err := s.source.ReadCards(ctx, filePath, limit)
	if err != nil {
		return err
	}

	for i := range cards {
		cards[i].Card.GameID = game.ID
	}

	if err := s.cards.CreateCards(ctx, cards); err != nil {
		return err
	}

	slog.Info("import finished", "game", gameCode, "cards", len(cards))

	return nil
}
