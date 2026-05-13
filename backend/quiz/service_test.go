package quiz_test

import (
	"context"
	"testing"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/quiz"
	"github.com/Cellul4r/go-quiz-app/backend/quiz/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetByID(t *testing.T) {
	t.Run("success public quiz", func(t *testing.T) {
		svc, quizRepo := setupService(t)

		requesterID := uuid.New()
		quizID := uuid.New()
		quizRepo.On("GetByID", mock.Anything, quizID).
			Return(domain.Quiz{
				ID:               quizID,
				AuthorID:         uuid.New(),
				Title:            "Test Quiz",
				Description:      "A quiz for testing",
				VisibilityStatus: domain.VisibilityPublic,
			}, nil)

		quiz, err := svc.GetByID(context.Background(), requesterID, quizID)

		require.NoError(t, err)
		assert.Equal(t, quizID, quiz.ID)
		assert.Equal(t, "Test Quiz", quiz.Title)
		assert.Equal(t, "A quiz for testing", quiz.Description)
		assert.Equal(t, domain.VisibilityPublic, quiz.VisibilityStatus)
		quizRepo.AssertExpectations(t)
	})

	t.Run("forbidden view private Quiz", func(t *testing.T) {
		svc, quizRepo := setupService(t)

		requesterID := uuid.New()
		quizID := uuid.New()
		AuthorID := uuid.New()
		quizRepo.On("GetByID", mock.Anything, quizID).
			Return(domain.Quiz{
				ID:               quizID,
				AuthorID:         AuthorID,
				Title:            "Test Quiz",
				Description:      "A quiz for testing",
				VisibilityStatus: domain.VisibilityPrivate,
			}, nil)

		_, err := svc.GetByID(context.Background(), requesterID, quizID)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		quizRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		svc, quizRepo := setupService(t)

		requesterID := uuid.New()
		quizID := uuid.New()
		quizRepo.On("GetByID", mock.Anything, quizID).
			Return(domain.Quiz{}, domain.ErrNotFound)

		_, err := svc.GetByID(context.Background(), requesterID, quizID)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		quizRepo.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc, quizRepo := setupService(t)

		requesterID := uuid.New()
		quizID := uuid.New()
		quizRepo.On("GetByID", mock.Anything, quizID).
			Return(domain.Quiz{}, domain.ErrInternalServerError)

		_, err := svc.GetByID(context.Background(), requesterID, quizID)

		assert.ErrorIs(t, err, domain.ErrInternalServerError)
		quizRepo.AssertExpectations(t)
	})
}

func setupService(t *testing.T) (*quiz.Service, *mocks.MockQuizRepository) {
	t.Helper()
	quizRepo := new(mocks.MockQuizRepository)
	svc := quiz.NewService(quizRepo, nil)

	return svc, quizRepo
}
