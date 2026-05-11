package mocks

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockQuizService struct {
	mock.Mock
}

func (m *MockQuizService) GetByID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) (domain.Quiz, error) {
	args := m.Called(ctx, requesterID, quizID)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizService) GetAll(ctx context.Context, requesterID *uuid.UUID, onlyPublic bool) ([]domain.Quiz, error) {
	args := m.Called(ctx, requesterID, onlyPublic)
	return args.Get(0).([]domain.Quiz), args.Error(1)
}

func (m *MockQuizService) Create(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error) {
	args := m.Called(ctx, requesterID, quiz)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizService) UpdateByID(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error) {
	args := m.Called(ctx, requesterID, quiz)
	return args.Get(0).(domain.Quiz), args.Error(1)
}

func (m *MockQuizService) DeleteByID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) error {
	args := m.Called(ctx, requesterID, quizID)
	return args.Error(0)
}
