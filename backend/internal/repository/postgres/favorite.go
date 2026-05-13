package postgres

import (
	"context"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewFavoriteRepository(db *gorm.DB, logger *slog.Logger) *FavoriteRepository {
	if logger == nil {
		logger = slog.Default()
	}
	return &FavoriteRepository{
		db:     db,
		logger: logger,
	}
}

func (m *FavoriteRepository) GetFavoritesByProfileID(ctx context.Context, profileID uuid.UUID) ([]domain.Quiz, error) {
	m.logger.Debug("favorite repository get favorites by profile ID query", "profile_id", profileID)
	var quizzes []domain.Quiz

	result := m.db.WithContext(ctx).Order("created_at desc").Joins("JOIN favorites f ON f.quiz_id = quizzes.id").Where("f.profile_id = ?", profileID).Find(&quizzes)
	if result.Error != nil {
		m.logger.Error("favorite repository get favorites by profile ID failed", "profile_id", profileID, "error", result.Error)
		return nil, domain.ErrInternalServerError
	}

	if quizzes == nil {
		quizzes = []domain.Quiz{}
	}
	return quizzes, nil
}

func (m *FavoriteRepository) AddFavorite(ctx context.Context, profileID uuid.UUID, quizID uuid.UUID) error {
	m.logger.Debug("favorite repository add favorite query", "profile_id", profileID, "quiz_id", quizID)
	result := m.db.WithContext(ctx).Create(&domain.Favorite{
		ProfileID: profileID,
		QuizID:    quizID,
	})
	if result.Error != nil {
		m.logger.Error("favorite repository add favorite failed", "profile_id", profileID, "quiz_id", quizID, "error", result.Error)
		return domain.ErrInternalServerError
	}
	return nil
}

func (m *FavoriteRepository) RemoveFavorite(ctx context.Context, profileID uuid.UUID, quizID uuid.UUID) error {
	m.logger.Debug("favorite repository remove favorite query", "profile_id", profileID, "quiz_id", quizID)
	result := m.db.WithContext(ctx).Delete(&domain.Favorite{}, "profile_id = ? AND quiz_id = ?", profileID, quizID)
	if result.Error != nil {
		m.logger.Error("favorite repository remove favorite failed", "profile_id", profileID, "quiz_id", quizID, "error", result.Error)
		return domain.ErrInternalServerError
	}
	return nil
}
