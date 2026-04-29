package quiz

import (
	"context"
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
)

type QuizRepository interface {
	GetByID(ctx context.Context, quizID uuid.UUID) (domain.Quiz, error)
	GetAll(ctx context.Context, filter domain.QuizFilter) ([]domain.Quiz, error)
	Create(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error)
	UpdateByID(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error)
	DeleteByID(ctx context.Context, quizID uuid.UUID) error
}

type Service struct {
	quizRepo QuizRepository
	logger   *slog.Logger
}

func NewService(quizRepo QuizRepository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{quizRepo: quizRepo, logger: logger}
}

func (s *Service) GetByID(ctx context.Context, requestID uuid.UUID, quizID uuid.UUID) (domain.Quiz, error) {

	quiz, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Error("quiz service get by id failed", "request_id", requestID.String(), "quiz_id", quizID.String(), "error", err)
		return domain.Quiz{}, err
	}

	if quiz.VisibilityStatus == domain.VisibilityPrivate &&
		quiz.ProfileID != requestID {
		s.logger.Warn("quiz service get by id unauthorized access", "request_id", requestID.String(), "quiz_id", quizID.String())
		return domain.Quiz{}, domain.ErrForbidden
	}
	return quiz, nil
}

func (s *Service) GetAll(ctx context.Context, requesterID *uuid.UUID, onlyPublic bool) ([]domain.Quiz, error) {

	if !onlyPublic && requesterID == nil {
		s.logger.Warn("quiz service get all missing requester id for non-public request")
		return nil, domain.ErrBadParamInput
	}

	filter := domain.QuizFilter{
		ProfileID:  requesterID,
		OnlyPublic: onlyPublic,
	}

	quizzes, err := s.quizRepo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("quiz service get all failed", "requester_id", requesterID, "error", err)
		return nil, err
	}

	return quizzes, nil
}

func (s *Service) Create(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error) {

	quiz.ID = uuid.New()
	quiz.ProfileID = requesterID

	created, err := s.quizRepo.Create(ctx, quiz)
	if err != nil {
		s.logger.Error("quiz service create failed", "request_id", requesterID.String(), "error", err)
		return domain.Quiz{}, err
	}

	return created, nil
}

func (s *Service) UpdateByID(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error) {
	existing, err := s.quizRepo.GetByID(ctx, quiz.ID)
	if err != nil {
		s.logger.Error("quiz service update get existing failed", "request_id", requesterID.String(), "quiz_id", quiz.ID.String(), "error", err)
		return domain.Quiz{}, err
	}

	if existing.ProfileID != requesterID {
		s.logger.Warn("quiz service update unauthorized access", "request_id", requesterID.String(), "quiz_id", quiz.ID.String())
		return domain.Quiz{}, domain.ErrForbidden
	}

	quiz.ProfileID = requesterID
	updatedQuiz, err := s.quizRepo.UpdateByID(ctx, quiz)
	if err != nil {
		s.logger.Error("quiz service update failed", "request_id", requesterID.String(), "quiz_id", quiz.ID.String(), "error", err)
		return domain.Quiz{}, err
	}

	return updatedQuiz, nil
}

func (s *Service) DeleteByID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) error {
	existing, err := s.quizRepo.GetByID(ctx, quizID)
	if err != nil {
		s.logger.Error("quiz service delete get existing failed", "request_id", requesterID.String(), "quiz_id", quizID.String(), "error", err)
		return err
	}

	if existing.ProfileID != requesterID {
		s.logger.Warn("quiz service delete unauthorized access", "request_id", requesterID.String(), "quiz_id", quizID.String())
		return domain.ErrForbidden
	}

	if err := s.quizRepo.DeleteByID(ctx, quizID); err != nil {
		s.logger.Error("quiz service delete failed", "request_id", requesterID.String(), "quiz_id", quizID.String(), "error", err)
		return err
	}

	return nil
}
