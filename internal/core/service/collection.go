package service

import (
	"context"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

type collectionService struct {
	repo port.CollectionRepository
}

func NewCollectionService(repo port.CollectionRepository) port.CollectionService {
	return &collectionService{repo: repo}
}

func (s *collectionService) CreateCollection(ctx context.Context, in domain.NewCollection) (domain.Collection, error) {
	return s.repo.CreateCollection(ctx, in)
}

func (s *collectionService) AddCardToCollection(ctx context.Context, in domain.CardToCollection) error {
	return s.repo.AddCardToCollection(ctx, in)
}

func (s *collectionService) ListCollections(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {
	return s.repo.ListCollections(ctx, userID)
}
