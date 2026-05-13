package mocks

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockQuizRepository struct {
	mock.Mock
}

func (m *MockQuizRepository) GetByID(ctx context.Context, quizID uuid.UUID) (domain.Quiz, error) {
	args := m.Called(ctx, quizID)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizRepository) GetAll(ctx context.Context, filter domain.QuizFilter) ([]domain.Quiz, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Quiz), args.Error(1)
}

func (m *MockQuizRepository) Create(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error) {
	args := m.Called(ctx, quiz)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizRepository) UpdateByID(ctx context.Context, quiz *domain.Quiz) (domain.Quiz, error) {
	args := m.Called(ctx, quiz)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizRepository) DeleteByID(ctx context.Context, quizID uuid.UUID) error {
	args := m.Called(ctx, quizID)
	return args.Error(0)
}
