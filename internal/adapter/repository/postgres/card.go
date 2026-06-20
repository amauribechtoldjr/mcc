package postgres

import (
	"context"
	"fmt"
	"strings"

	repo "github.com/amauribechtoldjr/mcc/internal/adapter/repository/postgres/sqlc"
	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
	for _, card := range cards {
		cardId, err := r.q.CreateCard(ctx, repo.CreateCardParams{
			OracleID: card.Card.OracleID,
			GameID:   card.Card.GameID,
		})

		if err != nil {
			fmt.Printf("CARD ERROR: %v \n", err)
			return mapError(err)
		}

		card.MTGCard.CardID = cardId

		cmcValue := decimal.NullDecimal{
			Decimal: decimal.NewFromFloat32(card.MTGCard.CMC),
			Valid:   true,
		}

		mtgCardParams := repo.CreateMTGCardParams{
			SetID:          card.MTGCard.SetID,
			Name:           card.MTGCard.Name,
			CardID:         card.MTGCard.CardID,
			Layout:         &card.MTGCard.Layout,
			Cmc:            cmcValue,
			ColorIdentity:  stringsToString(card.MTGCard.ColorIdentity),
			ColorIndicator: stringsToString(card.MTGCard.ColorIndicator),
			Colors:         stringsToString(card.MTGCard.Colors),
			ImgSmallUri:    &card.MTGCard.ImgSmallURI,
			ImgNormalUri:   &card.MTGCard.ImgNormalURI,
		}

		err = r.q.CreateMTGCard(ctx, mtgCardParams)
		if err != nil {
			fmt.Printf("MTG CARD ERROR: %v", err)
			return mapError(err)
		}
	}

	return nil
}

func stringsToString(strs []string) *string {
	var res string
	res = strings.Join(strs, ",")
	return &res
}

func toDomainCard(row repo.Card) domain.Card {
	return domain.Card{
		ID:       row.ID,
		OracleID: row.OracleID,
		GameID:   row.GameID,
	}
}
