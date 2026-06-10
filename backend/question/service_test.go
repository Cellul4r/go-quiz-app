package question

import (
	"context"
	"testing"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockQuestionRepository struct {
	mock.Mock
}

func (m *mockQuestionRepository) GetByID(ctx context.Context, questionID string) (*domain.Question, error) {
	args := m.Called(ctx, questionID)
	question, _ := args.Get(0).(*domain.Question)
	return question, args.Error(1)
}

func (m *mockQuestionRepository) GetByQuizID(ctx context.Context, quizID string) ([]domain.Question, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).([]domain.Question), args.Error(1)
}

func (m *mockQuestionRepository) Create(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	args := m.Called(ctx, question)
	created, _ := args.Get(0).(*domain.Question)
	return created, args.Error(1)
}

func (m *mockQuestionRepository) UpdateByID(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	args := m.Called(ctx, question)
	updated, _ := args.Get(0).(*domain.Question)
	return updated, args.Error(1)
}

func (m *mockQuestionRepository) DeleteByID(ctx context.Context, questionID string) error {
	args := m.Called(ctx, questionID)
	return args.Error(0)
}

type mockQuizRepository struct {
	mock.Mock
}

func (m *mockQuizRepository) GetByID(ctx context.Context, quizID uuid.UUID) (domain.Quiz, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func TestGetByID(t *testing.T) {
	t.Run("success public parent quiz", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		questionID := uuid.NewString()
		question := validQuestion(questionID, quizID)

		questionRepo.On("GetByID", mock.Anything, questionID).Return(question, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: uuid.New(), VisibilityStatus: domain.VisibilityPublic}, nil)

		got, err := svc.GetByID(context.Background(), requesterID, questionID)

		require.NoError(t, err)
		assert.Equal(t, questionID, got.ID)
		questionRepo.AssertExpectations(t)
		quizRepo.AssertExpectations(t)
	})

	t.Run("forbidden private parent quiz", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		questionID := uuid.NewString()
		question := validQuestion(questionID, quizID)

		questionRepo.On("GetByID", mock.Anything, questionID).Return(question, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: uuid.New(), VisibilityStatus: domain.VisibilityPrivate}, nil)

		_, err := svc.GetByID(context.Background(), requesterID, questionID)

		assert.ErrorIs(t, err, domain.ErrForbidden)
		questionRepo.AssertExpectations(t)
		quizRepo.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	t.Run("success when requester owns quiz", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		question := validQuestion("", quizID)

		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: requesterID, VisibilityStatus: domain.VisibilityPrivate}, nil)
		questionRepo.On("Create", mock.Anything, mock.MatchedBy(func(q *domain.Question) bool {
			_, err := uuid.Parse(q.ID)
			return err == nil && q.QuizID == quizID.String()
		})).Return(question, nil)

		created, err := svc.Create(context.Background(), requesterID, question)

		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		questionRepo.AssertExpectations(t)
		quizRepo.AssertExpectations(t)
	})

	t.Run("rejects invalid multiple choice options", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		question := validQuestion(uuid.NewString(), uuid.New())
		question.Options = domain.MultipleChoiceOptions{Choices: []domain.Choice{{ID: "a", Text: "A"}}, CorrectChoices: []string{"missing"}}

		_, err := svc.Create(context.Background(), uuid.New(), question)

		assert.ErrorIs(t, err, domain.ErrBadParamInput)
		questionRepo.AssertNotCalled(t, "Create")
		quizRepo.AssertNotCalled(t, "GetByID")
	})
}

func TestUpdateByID(t *testing.T) {
	t.Run("merges non-zero fields and preserves omitted fields", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		questionID := uuid.NewString()
		existing := validQuestion(questionID, quizID)
		existing.Content = "old content"
		existing.TimeLimitSeconds = 30
		existing.SortOrder = 4
		update := &domain.Question{ID: questionID, Content: "new content", SortOrder: 0}

		questionRepo.On("GetByID", mock.Anything, questionID).Return(existing, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: requesterID, VisibilityStatus: domain.VisibilityPrivate}, nil)
		questionRepo.On("UpdateByID", mock.Anything, mock.MatchedBy(func(q *domain.Question) bool {
			return q.ID == questionID && q.QuizID == quizID.String() && q.Content == "new content" && q.TimeLimitSeconds == 30 && q.SortOrder == 4
		})).Return(&domain.Question{ID: questionID, QuizID: quizID.String(), Content: "new content"}, nil)

		_, err := svc.UpdateByID(context.Background(), requesterID, update)

		require.NoError(t, err)
		questionRepo.AssertExpectations(t)
		quizRepo.AssertExpectations(t)
	})

	t.Run("forbidden when requester does not own quiz", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		quizID := uuid.New()
		questionID := uuid.NewString()
		existing := validQuestion(questionID, quizID)

		questionRepo.On("GetByID", mock.Anything, questionID).Return(existing, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: uuid.New(), VisibilityStatus: domain.VisibilityPrivate}, nil)

		_, err := svc.UpdateByID(context.Background(), uuid.New(), &domain.Question{ID: questionID, Content: "new"})

		assert.ErrorIs(t, err, domain.ErrForbidden)
		questionRepo.AssertNotCalled(t, "UpdateByID")
	})

	t.Run("rejects invalid updated type without matching options", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		questionID := uuid.NewString()
		existing := validQuestion(questionID, quizID)

		questionRepo.On("GetByID", mock.Anything, questionID).Return(existing, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: requesterID}, nil)

		_, err := svc.UpdateByID(context.Background(), requesterID, &domain.Question{ID: questionID, QuestionType: domain.TypedAnswer})

		assert.ErrorIs(t, err, domain.ErrBadParamInput)
		questionRepo.AssertNotCalled(t, "UpdateByID")
	})
}

func TestDeleteByID(t *testing.T) {
	t.Run("success when requester owns quiz", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)
		requesterID := uuid.New()
		quizID := uuid.New()
		questionID := uuid.NewString()
		existing := validQuestion(questionID, quizID)

		questionRepo.On("GetByID", mock.Anything, questionID).Return(existing, nil)
		quizRepo.On("GetByID", mock.Anything, quizID).Return(domain.Quiz{ID: quizID, AuthorID: requesterID}, nil)
		questionRepo.On("DeleteByID", mock.Anything, questionID).Return(nil)

		err := svc.DeleteByID(context.Background(), requesterID, questionID)

		require.NoError(t, err)
		questionRepo.AssertExpectations(t)
		quizRepo.AssertExpectations(t)
	})

	t.Run("rejects invalid id", func(t *testing.T) {
		svc, questionRepo, quizRepo := setupService(t)

		err := svc.DeleteByID(context.Background(), uuid.New(), "bad-id")

		assert.ErrorIs(t, err, domain.ErrBadParamInput)
		questionRepo.AssertNotCalled(t, "GetByID")
		quizRepo.AssertNotCalled(t, "GetByID")
	})
}

func setupService(t *testing.T) (*Service, *mockQuestionRepository, *mockQuizRepository) {
	t.Helper()
	questionRepo := new(mockQuestionRepository)
	quizRepo := new(mockQuizRepository)
	return NewService(questionRepo, quizRepo, nil), questionRepo, quizRepo
}

func validQuestion(questionID string, quizID uuid.UUID) *domain.Question {
	return &domain.Question{
		ID:               questionID,
		QuizID:           quizID.String(),
		QuestionType:     domain.MultipleChoice,
		Content:          "question content",
		TimeLimitSeconds: 15,
		SortOrder:        1,
		Options: domain.MultipleChoiceOptions{
			Choices: []domain.Choice{
				{ID: "a", Text: "A"},
				{ID: "b", Text: "B"},
			},
			CorrectChoices: []string{"a"},
		},
	}
}
