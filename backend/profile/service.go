package profile

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type ProfileRepository interface {
	GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error)
	GetByUsername(ctx context.Context, username string) (domain.Profile, error)
	UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error)
}

type Service struct {
	profileRepo ProfileRepository
}

func NewService(profileRepo ProfileRepository) *Service {
	return &Service{profileRepo: profileRepo}
}

func (s *Service) GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	profile, err := s.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		return domain.Profile{}, err
	}
	return profile, nil
}

func (s *Service) UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error) {
	// validation
	existing, err := s.profileRepo.GetByUsername(ctx, profile.Username)
	if err == nil && existing.ID != profile.ID {
		return domain.Profile{}, domain.ErrUserNameConflict
	}

	return s.profileRepo.UpdateByID(ctx, profile)
}
