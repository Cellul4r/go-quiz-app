package mocks

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFavoriteService struct {
	mock.Mock
}

func (m *MockFavoriteService) GetFavoritesByProfileID(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID) ([]domain.Quiz, error) {
	args := m.Called(ctx, requesterID, profileID)
	return args.Get(0).([]domain.Quiz), args.Error(1)
}

func (m *MockFavoriteService) AddFavorite(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID, quizID uuid.UUID) error {
	args := m.Called(ctx, requesterID, profileID, quizID)
	return args.Error(0)
}

func (m *MockFavoriteService) RemoveFavorite(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID, quizID uuid.UUID) error {
	args := m.Called(ctx, requesterID, profileID, quizID)
	return args.Error(0)
}
