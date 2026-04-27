package profile

import (
	"context"
	"errors"
	"log/slog"

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
	logger      *slog.Logger
}

func NewService(profileRepo ProfileRepository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{profileRepo: profileRepo, logger: logger}
}

func (s *Service) GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	s.logger.Debug("profile service get by id start", "profile_id", profileID.String())
	profile, err := s.profileRepo.GetByID(ctx, profileID)
	if err != nil {
		s.logger.Error("profile service get by id failed", "profile_id", profileID.String(), "error", err)
		return domain.Profile{}, err
	}
	return profile, nil
}

func (s *Service) UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error) {
	s.logger.Debug("profile service update start", "profile_id", profile.ID.String(), "username", profile.Username)

	existing, err := s.profileRepo.GetByID(ctx, profile.ID)
	if err != nil {
		s.logger.Error("profile service update get existing failed", "profile_id", profile.ID.String(), "error", err)
		return domain.Profile{}, err
	}

	if profile.Username != "" {
		existing, err = s.profileRepo.GetByUsername(ctx, profile.Username)
		if err == nil && existing.ID != profile.ID {
			s.logger.Warn("profile service update username conflict", "profile_id", profile.ID.String(), "username", profile.Username)
			return domain.Profile{}, domain.ErrUserNameConflict
		} else if err != nil && !errors.Is(err, domain.ErrNotFound) {
			s.logger.Error("profile service update get by username failed", "username", profile.Username, "error", err)
			return domain.Profile{}, domain.ErrInternalServerError
		}
	}

	updated, updateErr := s.profileRepo.UpdateByID(ctx, profile)
	if updateErr != nil {
		s.logger.Error("profile service update failed", "profile_id", profile.ID.String(), "error", updateErr)
		return domain.Profile{}, updateErr
	}

	return updated, nil
}
