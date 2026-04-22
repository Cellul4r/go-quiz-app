package postgres

import (
	"context"
	"errors"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (m *ProfileRepository) GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	var profile domain.Profile
	if result := m.db.WithContext(ctx).First(&profile, "id = ?", profileID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Profile{}, domain.ErrNotFound
		}
		return domain.Profile{}, domain.ErrInternalServerError
	}
	return profile, nil
}

func (m *ProfileRepository) GetByUsername(ctx context.Context, username string) (domain.Profile, error) {
	var profile domain.Profile
	if result := m.db.WithContext(ctx).First(&profile, "username = ?", username); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.Profile{}, domain.ErrNotFound
		}
		return domain.Profile{}, domain.ErrInternalServerError
	}
	return profile, nil
}

func (m *ProfileRepository) UpdateByID(ctx context.Context, profile *domain.Profile) error {
	result := m.db.WithContext(ctx).
		Model(&domain.Profile{}).
		Where("id = ?", profile.ID).
		Updates(map[string]interface{}{
			"username":   profile.Username,
			"full_name":  profile.FullName,
			"avatar_url": profile.AvatarURL,
		})

	if result.Error != nil {
		return domain.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
