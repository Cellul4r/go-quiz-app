package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres/model"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewQuestionRepository(db *gorm.DB, logger *slog.Logger) *QuestionRepository {
	if logger == nil {
		logger = slog.Default()
	}
	return &QuestionRepository{db: db, logger: logger}
}

func (m *QuestionRepository) GetByID(ctx context.Context, questionID string) (*domain.Question, error) {
	m.logger.Debug("question repository get by id query", "question_id", questionID)
	var questionModel model.QuestionModel
	if result := m.db.WithContext(ctx).First(&questionModel, "id = ?", questionID); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			m.logger.Warn("question repository get by id not found", "question_id", questionID)
			return nil, domain.ErrNotFound
		}
		m.logger.Error("question repository get by id failed", "question_id", questionID, "error", result.Error)
		return nil, domain.ErrInternalServerError
	}

	return model.ToDomainQuestion(&questionModel)
}

func (m *QuestionRepository) GetByQuizID(ctx context.Context, quizID string) ([]domain.Question, error) {
	m.logger.Debug("question repository get by quiz id query", "quiz_id", quizID)
	var questions []model.QuestionModel
	if result := m.db.WithContext(ctx).
		Where("quiz_id = ?", quizID).
		Order("sort_order asc").
		Find(&questions); result.Error != nil {
		m.logger.Error("question repository get by quiz id failed", "quiz_id", quizID, "error", result.Error)
		return nil, domain.ErrInternalServerError
	}

	if questions == nil {
		questions = []model.QuestionModel{}
	}

	domainQuestions := make([]domain.Question, len(questions))
	for i, q := range questions {
		domainQ, err := model.ToDomainQuestion(&q)
		if err != nil {
			m.logger.Error("question repository get by quiz id failed to convert to domain model", "quiz_id", quizID, "error", err)
			return nil, domain.ErrInternalServerError
		}
		domainQuestions[i] = *domainQ
	}
	return domainQuestions, nil
}

func (m *QuestionRepository) Create(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	m.logger.Debug("question repository create query", "question_id", question.ID, "quiz_id", question.QuizID)
	questionModel, err := model.ToModelQuestion(question)
	if err != nil {
		m.logger.Error("question repository create failed to convert to model", "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return nil, domain.ErrInternalServerError
	}
	if err := m.db.WithContext(ctx).Create(&questionModel).Error; err != nil {
		m.logger.Error("question repository create failed", "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return nil, domain.ErrInternalServerError
	}
	return model.ToDomainQuestion(questionModel)
}

func (m *QuestionRepository) UpdateByID(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	m.logger.Debug("question repository update query", "question_id", question.ID, "quiz_id", question.QuizID)
	questionModel, err := model.ToModelQuestion(question)
	if err != nil {
		m.logger.Error("question repository update failed to convert to model", "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return nil, domain.ErrInternalServerError
	}
	if err := m.db.WithContext(ctx).Where("id = ?", question.ID).Updates(&questionModel).Error; err != nil {
		m.logger.Error("question repository update failed", "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return nil, domain.ErrInternalServerError
	}
	return model.ToDomainQuestion(questionModel)
}

func (m *QuestionRepository) DeleteByID(ctx context.Context, questionID string) error {
	m.logger.Debug("question repository delete query", "question_id", questionID)
	if err := m.db.WithContext(ctx).Where("id = ?", questionID).Delete(&model.QuestionModel{}).Error; err != nil {
		m.logger.Error("question repository delete failed", "question_id", questionID, "error", err)
		return domain.ErrInternalServerError
	}
	return nil
}
