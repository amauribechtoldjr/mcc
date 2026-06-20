package port

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/google/uuid"
)

type MTGSetRepository interface {
	CreateMTGSet(ctx context.Context, set domain.MTGSet) (uuid.UUID, error)
}
