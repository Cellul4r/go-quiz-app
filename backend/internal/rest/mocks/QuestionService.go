package mocks

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockQuestionService struct {
	mock.Mock
}

func (m *MockQuestionService) GetByID(ctx context.Context, requesterID uuid.UUID, questionID string) (*domain.Question, error) {
	args := m.Called(ctx, requesterID, questionID)
	question, _ := args.Get(0).(*domain.Question)
	return question, args.Error(1)
}

func (m *MockQuestionService) GetByQuizID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) ([]domain.Question, error) {
	args := m.Called(ctx, requesterID, quizID)
	return args.Get(0).([]domain.Question), args.Error(1)
}

func (m *MockQuestionService) Create(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error) {
	args := m.Called(ctx, requesterID, question)
	created, _ := args.Get(0).(*domain.Question)
	return created, args.Error(1)
}

func (m *MockQuestionService) UpdateByID(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error) {
	args := m.Called(ctx, requesterID, question)
	updated, _ := args.Get(0).(*domain.Question)
	return updated, args.Error(1)
}

func (m *MockQuestionService) DeleteByID(ctx context.Context, requesterID uuid.UUID, questionID string) error {
	args := m.Called(ctx, requesterID, questionID)
	return args.Error(0)
}
