package port

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
)

type GameRepository interface {
	FindGameByCode(ctx context.Context, code string) (domain.Game, error)
}
