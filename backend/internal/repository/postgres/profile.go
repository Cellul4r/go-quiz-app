package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewProfileRepository(db *gorm.DB, logger *slog.Logger) *ProfileRepository {
	if logger == nil {
		logger = slog.Default()
	}

	return &ProfileRepository{db: db, logger: logger}
}

func (m *ProfileRepository) GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	m.logger.Debug("profile repository get by id query", "profile_id", profileID.String())
	var profile domain.Profile
	if result := m.db.WithContext(ctx).First(&profile, "id = ?", profileID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			m.logger.Warn("profile repository get by id not found", "profile_id", profileID.String())
			return domain.Profile{}, domain.ErrNotFound
		}
		m.logger.Error("profile repository get by id failed", "profile_id", profileID.String(), "error", result.Error)
		return domain.Profile{}, domain.ErrInternalServerError
	}
	return profile, nil
}

func (m *ProfileRepository) GetByUsername(ctx context.Context, username string) (domain.Profile, error) {
	m.logger.Debug("profile repository get by username query", "username", username)
	var profile domain.Profile
	if result := m.db.WithContext(ctx).First(&profile, "username = ?", username); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			m.logger.Warn("profile repository get by username not found", "username", username)
			return domain.Profile{}, domain.ErrNotFound
		}
		m.logger.Error("profile repository get by username failed", "username", username, "error", result.Error)
		return domain.Profile{}, domain.ErrInternalServerError
	}
	return profile, nil
}

func (m *ProfileRepository) UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error) {
	m.logger.Debug("profile repository update query", "profile_id", profile.ID.String(), "profile", profile)
	result := m.db.WithContext(ctx).
		Model(&profile).
		Where("id = ?", profile.ID).
		Updates(map[string]interface{}{
			"username":   profile.Username,
			"full_name":  profile.FullName,
			"avatar_url": profile.AvatarURL,
		})

	if result.Error != nil {
		m.logger.Error("profile repository update failed", "profile_id", profile.ID.String(), "error", result.Error)
		return domain.Profile{}, domain.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		m.logger.Warn("profile repository update not found", "profile_id", profile.ID.String())
		return domain.Profile{}, domain.ErrNotFound
	}

	updated, err := m.GetByID(ctx, profile.ID)
	if err != nil {
		m.logger.Error("profile repository update fetch updated failed", "profile_id", profile.ID.String(), "error", err)
		return domain.Profile{}, err
	}
	return updated, nil
}
