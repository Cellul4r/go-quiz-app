package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewQuizRepository(db *gorm.DB, logger *slog.Logger) *QuizRepository {
	if logger == nil {
		logger = slog.Default()
	}
	return &QuizRepository{db: db, logger: logger}
}

func (m *QuizRepository) GetByID(ctx context.Context, quizID uuid.UUID) (domain.Quiz, error) {
	m.logger.Debug("quiz repository get by id query", "quiz_id", quizID.String())
	var quiz domain.Quiz
	if result := m.db.WithContext(ctx).First(&quiz, "id = ?", quizID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			m.logger.Warn("quiz repository get by id not found", "quiz_id", quizID.String())
			return domain.Quiz{}, domain.ErrNotFound
		}
		m.logger.Error("quiz repository get by id failed", "quiz_id", quizID.String(), "error", result.Error)
		return domain.Quiz{}, domain.ErrInternalServerError
	}
	return quiz, nil
}

func (m *QuizRepository) GetAll(ctx context.Context, filter domain.QuizFilter) ([]domain.Quiz, error) {
	m.logger.Debug("quiz repository get all query", "filter", filter)

	query := m.db.WithContext(ctx)
	if filter.OnlyPublic {
		query = query.Where("visibility_status = ?", domain.VisibilityPublic)
	}

	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}

	var quizzes []domain.Quiz
	if result := query.Find(&quizzes); result.Error != nil {
		m.logger.Error("quiz repository get all failed", "filter", filter, "error", result.Error)
		return nil, domain.ErrInternalServerError
	}

	if quizzes == nil {
		quizzes = []domain.Quiz{}
	}
	return quizzes, nil
}

func (m *QuizRepository) Create(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error) {
	m.logger.Debug("quiz repository create query", "quiz", quiz)
	if result := m.db.WithContext(ctx).Create(quiz); result.Error != nil {
		m.logger.Error("quiz repository create failed", "quiz", quiz, "error", result.Error)
		return domain.Quiz{}, domain.ErrInternalServerError
	}
	return *quiz, nil
}

func (m *QuizRepository) UpdateByID(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error) {
	m.logger.Debug("quiz repository update query", "quiz_id", quiz.ID.String(), "quiz", quiz)

	updates := map[string]any{
		"description": quiz.Description,
	}

	if quiz.Title != "" {
		updates["title"] = quiz.Title
	}
	if quiz.VisibilityStatus != "" {
		updates["visibility_status"] = quiz.VisibilityStatus
	}

	result := m.db.WithContext(ctx).
		Model(&quiz).
		Updates(updates)

	if result.Error != nil {
		m.logger.Error("quiz repository update failed", "quiz_id", quiz.ID.String(), "error", result.Error)
		return domain.Quiz{}, domain.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		m.logger.Warn("quiz repository update not found", "quiz_id", quiz.ID.String())
		return domain.Quiz{}, domain.ErrNotFound
	}

	updated, err := m.GetByID(ctx, quiz.ID)
	if err != nil {
		m.logger.Error("quiz repository update fetch updated failed", "quiz_id", quiz.ID.String(), "error", err)
		return domain.Quiz{}, err
	}
	return updated, nil
}

func (m *QuizRepository) DeleteByID(ctx context.Context, quizID uuid.UUID) error {
	m.logger.Debug("quiz repository delete query", "quiz_id", quizID.String())
	result := m.db.WithContext(ctx).Delete(&domain.Quiz{}, "id = ?", quizID)

	if result.Error != nil {
		m.logger.Error("quiz repository delete failed", "quiz_id", quizID.String(), "error", result.Error)
		return domain.ErrInternalServerError
	}
	if result.RowsAffected == 0 {
		m.logger.Warn("quiz repository delete not found", "quiz_id", quizID.String())
		return domain.ErrNotFound
	}
	return nil
}
