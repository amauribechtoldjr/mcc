package collections

import (
	"context"

	repo "github.com/amauribechtoldjr/mcc/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) CollectionsService {
	return &svc{repo: repo}
}

func (s *svc) CreateCollection(ctx context.Context, collectionData repo.CreateCollectionParams) (repo.Collection, error) {
	return s.repo.CreateCollection(ctx, collectionData)
}
