package favorite

import (
	"context"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type FavoriteRepository interface {
	GetFavoritesByProfileID(ctx context.Context, profileID uuid.UUID) ([]domain.Quiz, error)
	AddFavorite(ctx context.Context, profileID uuid.UUID, quizID uuid.UUID) error
	RemoveFavorite(ctx context.Context, profileID uuid.UUID, quizID uuid.UUID) error
}

type Service struct {
	favoriteRepo FavoriteRepository
	logger       *slog.Logger
}

func NewService(favoriteRepo FavoriteRepository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		favoriteRepo: favoriteRepo,
		logger:       logger,
	}
}

func (s *Service) GetFavoritesByProfileID(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID) ([]domain.Quiz, error) {
	if requesterID != profileID {
		s.logger.Warn("requester ID does not match profile ID, returning empty favorites", "requesterID", requesterID, "profileID", profileID)
		return []domain.Quiz{}, nil
	}
	return s.favoriteRepo.GetFavoritesByProfileID(ctx, profileID)
}

func (s *Service) AddFavorite(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID, quizID uuid.UUID) error {
	if requesterID != profileID {
		s.logger.Warn("requester ID does not match profile ID, cannot add favorite", "requesterID", requesterID, "profileID", profileID)
		return nil
	}

	s.favoriteRepo.AddFavorite(ctx, profileID, quizID)
	return nil
}

func (s *Service) RemoveFavorite(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID, quizID uuid.UUID) error {
	if requesterID != profileID {
		s.logger.Warn("requester ID does not match profile ID, cannot remove favorite", "requesterID", requesterID, "profileID", profileID)
		return nil
	}

	s.favoriteRepo.RemoveFavorite(ctx, profileID, quizID)
	return nil
}
