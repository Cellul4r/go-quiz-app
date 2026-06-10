package question

import (
	"context"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type QuestionRepository interface {
	GetByID(ctx context.Context, questionID string) (*domain.Question, error)
	GetByQuizID(ctx context.Context, quizID string) ([]domain.Question, error)
	Create(ctx context.Context, question *domain.Question) (*domain.Question, error)
	UpdateByID(ctx context.Context, question *domain.Question) (*domain.Question, error)
	DeleteByID(ctx context.Context, questionID string) error
}

type QuizRepository interface {
	GetByID(ctx context.Context, quizID uuid.UUID) (domain.Quiz, error)
}

type Service struct {
	questionRepo QuestionRepository
	quizRepo     QuizRepository
	logger       *slog.Logger
}

func NewService(questionRepo QuestionRepository, quizRepo QuizRepository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{questionRepo: questionRepo, quizRepo: quizRepo, logger: logger}
}

func (s *Service) GetByID(ctx context.Context, requesterID uuid.UUID, questionID string) (*domain.Question, error) {
	question, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		s.logger.Error("question service get by id failed", "request_id", requesterID.String(), "question_id", questionID, "error", err)
		return nil, err
	}

	if err := s.canAccessQuestion(ctx, requesterID, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *Service) GetByQuizID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) ([]domain.Question, error) {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Error("question service get by quiz id get parent failed", "request_id", requesterID.String(), "quiz_id", quizID.String(), "error", err)
		return nil, err
	}
	if err := s.canAccessQuiz(requesterID, quiz); err != nil {
		return nil, err
	}

	questions, err := s.questionRepo.GetByQuizID(ctx, quizID.String())
	if err != nil {
		s.logger.Error("question service get by quiz id failed", "request_id", requesterID.String(), "quiz_id", quizID.String(), "error", err)
		return nil, err
	}
	return questions, nil
}

func (s *Service) Create(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error) {
	if question == nil {
		return nil, domain.ErrBadParamInput
	}
	if question.ID == "" {
		question.ID = uuid.NewString()
	}
	if err := validateQuestion(question, true); err != nil {
		return nil, err
	}

	quizID, err := uuid.Parse(question.QuizID)
	if err != nil {
		return nil, domain.ErrBadParamInput
	}
	if err := s.canMutateQuiz(ctx, requesterID, quizID); err != nil {
		return nil, err
	}

	created, err := s.questionRepo.Create(ctx, question)
	if err != nil {
		s.logger.Error("question service create failed", "request_id", requesterID.String(), "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return nil, err
	}
	return created, nil
}

func (s *Service) UpdateByID(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error) {
	if question == nil || question.ID == "" {
		return nil, domain.ErrBadParamInput
	}
	if _, err := uuid.Parse(question.ID); err != nil {
		return nil, domain.ErrBadParamInput
	}

	existing, err := s.questionRepo.GetByID(ctx, question.ID)
	if err != nil {
		s.logger.Error("question service update get existing failed", "request_id", requesterID.String(), "question_id", question.ID, "error", err)
		return nil, err
	}

	quizID, err := uuid.Parse(existing.QuizID)
	if err != nil {
		return nil, domain.ErrBadParamInput
	}
	if err := s.canMutateQuiz(ctx, requesterID, quizID); err != nil {
		return nil, err
	}

	merged := mergeQuestionUpdate(existing, question)
	if err := validateQuestion(merged, true); err != nil {
		return nil, err
	}

	updated, err := s.questionRepo.UpdateByID(ctx, merged)
	if err != nil {
		s.logger.Error("question service update failed", "request_id", requesterID.String(), "question_id", question.ID, "error", err)
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeleteByID(ctx context.Context, requesterID uuid.UUID, questionID string) error {
	if _, err := uuid.Parse(questionID); err != nil {
		return domain.ErrBadParamInput
	}

	existing, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		s.logger.Error("question service delete get existing failed", "request_id", requesterID.String(), "question_id", questionID, "error", err)
		return err
	}

	quizID, err := uuid.Parse(existing.QuizID)
	if err != nil {
		return domain.ErrBadParamInput
	}
	if err := s.canMutateQuiz(ctx, requesterID, quizID); err != nil {
		return err
	}

	if err := s.questionRepo.DeleteByID(ctx, questionID); err != nil {
		s.logger.Error("question service delete failed", "request_id", requesterID.String(), "question_id", questionID, "error", err)
		return err
	}
	return nil
}

func (s *Service) canAccessQuestion(ctx context.Context, requesterID uuid.UUID, question *domain.Question) error {
	if question == nil {
		return domain.ErrNotFound
	}
	quizID, err := uuid.Parse(question.QuizID)
	if err != nil {
		return domain.ErrBadParamInput
	}
	quiz, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Error("question service get parent quiz failed", "request_id", requesterID.String(), "question_id", question.ID, "quiz_id", question.QuizID, "error", err)
		return err
	}
	return s.canAccessQuiz(requesterID, quiz)
}

func (s *Service) canAccessQuiz(requesterID uuid.UUID, quiz domain.Quiz) error {
	if quiz.VisibilityStatus == domain.VisibilityPrivate && quiz.AuthorID != requesterID {
		s.logger.Warn("question service unauthorized private quiz access", "request_id", requesterID.String(), "quiz_id", quiz.ID.String())
		return domain.ErrForbidden
	}
	return nil
}

func (s *Service) canMutateQuiz(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) error {
	quiz, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Error("question service get parent quiz for mutation failed", "request_id", requesterID.String(), "quiz_id", quizID.String(), "error", err)
		return err
	}
	if quiz.AuthorID != requesterID {
		s.logger.Warn("question service unauthorized mutation", "request_id", requesterID.String(), "quiz_id", quizID.String())
		return domain.ErrForbidden
	}
	return nil
}

func mergeQuestionUpdate(existing *domain.Question, update *domain.Question) *domain.Question {
	merged := *existing
	if update.Content != "" {
		merged.Content = update.Content
	}
	if update.ImageURL != "" {
		merged.ImageURL = update.ImageURL
	}
	if update.TimeLimitSeconds != 0 {
		merged.TimeLimitSeconds = update.TimeLimitSeconds
	}
	if update.SortOrder != 0 {
		merged.SortOrder = update.SortOrder
	}
	if update.QuestionType != "" {
		merged.QuestionType = update.QuestionType
	}
	if update.Options != nil {
		merged.Options = update.Options
	}
	return &merged
}

func validateQuestion(question *domain.Question, requireContent bool) error {
	if question == nil {
		return domain.ErrBadParamInput
	}
	if question.ID != "" {
		if _, err := uuid.Parse(question.ID); err != nil {
			return domain.ErrBadParamInput
		}
	}
	if _, err := uuid.Parse(question.QuizID); err != nil {
		return domain.ErrBadParamInput
	}
	if requireContent && question.Content == "" {
		return domain.ErrBadParamInput
	}
	if question.TimeLimitSeconds < 0 || question.SortOrder < 0 {
		return domain.ErrBadParamInput
	}
	if !isValidQuestionType(question.QuestionType) {
		return domain.ErrBadParamInput
	}
	return validateOptions(question.QuestionType, question.Options)
}

func isValidQuestionType(questionType domain.QuestionType) bool {
	switch questionType {
	case domain.MultipleChoice, domain.TrueFalse, domain.TypedAnswer:
		return true
	default:
		return false
	}
}

func validateOptions(questionType domain.QuestionType, options any) error {
	switch questionType {
	case domain.MultipleChoice:
		opts, ok := options.(domain.MultipleChoiceOptions)
		if !ok {
			return domain.ErrBadParamInput
		}
		if len(opts.Choices) < 2 || len(opts.CorrectChoices) == 0 {
			return domain.ErrBadParamInput
		}
		choiceIDs := make(map[string]struct{}, len(opts.Choices))
		for _, choice := range opts.Choices {
			if choice.ID == "" || choice.Text == "" {
				return domain.ErrBadParamInput
			}
			if _, exists := choiceIDs[choice.ID]; exists {
				return domain.ErrBadParamInput
			}
			choiceIDs[choice.ID] = struct{}{}
		}
		for _, correctChoice := range opts.CorrectChoices {
			if _, exists := choiceIDs[correctChoice]; !exists {
				return domain.ErrBadParamInput
			}
		}
		return nil
	case domain.TrueFalse:
		if _, ok := options.(domain.TrueFalseOptions); !ok {
			return domain.ErrBadParamInput
		}
		return nil
	case domain.TypedAnswer:
		opts, ok := options.(domain.TypedAnswerOptions)
		if !ok || opts.CorrectAnswer == "" {
			return domain.ErrBadParamInput
		}
		return nil
	default:
		return domain.ErrBadParamInput
	}
}
