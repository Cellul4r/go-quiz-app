package mocks

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Profile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) GetByUsername(ctx context.Context, username string) (domain.Profile, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.Profile), args.Error(1)
}

func (m *MockProfileRepository) UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error) {
	args := m.Called(ctx, profile)
	return args.Get(0).(domain.Profile), args.Error(1)
}
