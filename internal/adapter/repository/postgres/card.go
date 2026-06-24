package postgres

import (
	"context"
	"strings"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const batchSize = 1000

type cardRepository struct {
	q repo.Querier
}

func NewCardRepository(q repo.Querier) port.CardRepository {
	return &cardRepository{q: q}
}

func (r *cardRepository) ListCards(ctx context.Context) ([]domain.Card, error) {
	rows, err := r.q.ListCards(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	cards := make([]domain.Card, 0, len(rows))
	for _, row := range rows {
		cards = append(cards, toDomainCard(row))
	}

	return cards, nil
}

func (r *cardRepository) FindCardByID(ctx context.Context, id uuid.UUID) (domain.Card, error) {
	row, err := r.q.FindCardById(ctx, id)
	if err != nil {
		return domain.Card{}, mapError(err)
	}

	return toDomainCard(row), nil
}

func (r *cardRepository) CreateCards(ctx context.Context, cards []domain.ImportCard) error {
	params := make([]repo.UpsertMTGCardParams, 0, len(cards))
	for _, card := range cards {
		cmcValue := decimal.NullDecimal{
			Decimal: decimal.NewFromFloat32(card.MTGCard.CMC),
			Valid:   true,
		}

		params = append(params, repo.UpsertMTGCardParams{
			GameID:          card.Card.GameID,
			SetID:           card.MTGCard.SetID,
			OracleID:        card.MTGCard.OracleID,
			Lang:            card.MTGCard.Lang,
			CollectorNumber: card.MTGCard.CollectorNumber,
			Name:            card.MTGCard.Name,
			PrintedTypeLine: nilIfEmpty(card.MTGCard.PrintedTypeLine),
			PrintedText:     nilIfEmpty(card.MTGCard.PrintedText),
			FlavorText:      nilIfEmpty(card.MTGCard.FlavorText),
			Layout:          nilIfEmpty(card.MTGCard.Layout),
			Cmc:             cmcValue,
			ColorIdentity:   stringsToString(card.MTGCard.ColorIdentity),
			ColorIndicator:  stringsToString(card.MTGCard.ColorIndicator),
			Colors:          stringsToString(card.MTGCard.Colors),
			ImgSmallUri:     nilIfEmpty(card.MTGCard.ImgSmallURI),
			ImgNormalUri:    nilIfEmpty(card.MTGCard.ImgNormalURI),
			LastImportID:    card.MTGCard.LastImportId,
		})
	}

	for start := 0; start < len(params); start += batchSize {
		end := min(start+batchSize, len(params))

		if err := r.upsertBatch(ctx, params[start:end]); err != nil {
			return err
		}
	}

	return nil
}

func (r *cardRepository) upsertBatch(ctx context.Context, params []repo.UpsertMTGCardParams) error {
	batch := r.q.UpsertMTGCard(ctx, params)
	defer batch.Close()

	var firstErr error
	batch.Exec(func(_ int, err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	})
	if firstErr != nil {
		return mapError(firstErr)
	}

	return nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func stringsToString(strs []string) *string {
	var res string
	res = strings.Join(strs, ",")
	return &res
}

func toDomainCard(row repo.Card) domain.Card {
	return domain.Card{
		ID:     row.ID,
		GameID: row.GameID,
	}
}
