package profile_test

import (
	"context"
	"testing"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/profile"
	"github.com/Cellul4r/go-quiz-app/backend/profile/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{
				ID:        profileID,
				Username:  "testuser",
				FullName:  "test user",
				AvatarURL: "http://example.com/avatar.png",
			}, nil)

		profile, err := svc.GetByID(context.Background(), profileID)

		require.NoError(t, err)
		assert.Equal(t, profileID, profile.ID)
		assert.Equal(t, "testuser", profile.Username)
		assert.Equal(t, "test user", profile.FullName)
		assert.Equal(t, "http://example.com/avatar.png", profile.AvatarURL)

		profileRepo.AssertExpectations(t)
	})

	t.Run("success with empty username", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{
				ID:        profileID,
				Username:  "",
				FullName:  "test user",
				AvatarURL: "http://example.com/avatar.png",
			}, nil)

		profile, err := svc.GetByID(context.Background(), profileID)

		require.NoError(t, err)
		assert.Equal(t, profileID, profile.ID)
		assert.Equal(t, "", profile.Username)
		assert.Equal(t, "test user", profile.FullName)
		assert.Equal(t, "http://example.com/avatar.png", profile.AvatarURL)

		profileRepo.AssertNotCalled(t, "GetByUsername")
		profileRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrNotFound)

		_, err := svc.GetByID(context.Background(), profileID)

		require.ErrorIs(t, err, domain.ErrNotFound)
		profileRepo.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrInternalServerError)

		_, err := svc.GetByID(context.Background(), profileID)

		require.ErrorIs(t, err, domain.ErrInternalServerError)
		profileRepo.AssertExpectations(t)
	})
}

func TestUpdateByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:        profileID,
			Username:  "updateduser",
			FullName:  "updated user",
			AvatarURL: "http://example.com/updated-avatar.png",
		}
		existing := domain.Profile{
			ID:        profileID,
			Username:  "testuser",
			FullName:  "test user",
			AvatarURL: "http://example.com/avatar.png",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("GetByUsername", mock.Anything, "updateduser").
			Return(domain.Profile{}, domain.ErrNotFound)
		profileRepo.On("UpdateByID", mock.Anything, input).
			Return(*input, nil)

		updated, err := svc.UpdateByID(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, input.ID, updated.ID)
		assert.Equal(t, input.Username, updated.Username)
		assert.Equal(t, input.FullName, updated.FullName)
		assert.Equal(t, input.AvatarURL, updated.AvatarURL)

		profileRepo.AssertExpectations(t)
	})

	t.Run("success with empty username", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:        profileID,
			Username:  "",
			FullName:  "updated user",
			AvatarURL: "http://example.com/updated-avatar.png",
		}
		existing := domain.Profile{
			ID:        profileID,
			Username:  "testuser",
			FullName:  "test user",
			AvatarURL: "http://example.com/avatar.png",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("UpdateByID", mock.Anything, input).
			Return(*input, nil)

		updated, err := svc.UpdateByID(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, input.ID, updated.ID)
		assert.Equal(t, "", updated.Username)
		assert.Equal(t, input.FullName, updated.FullName)
		assert.Equal(t, input.AvatarURL, updated.AvatarURL)

		profileRepo.AssertNotCalled(t, "GetByUsername")
		profileRepo.AssertExpectations(t)
	})

	t.Run("username conflict", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "existinguser",
		}
		existing := domain.Profile{
			ID:       profileID,
			Username: "testuser",
		}
		conflict := domain.Profile{
			ID:       uuid.New(),
			Username: "existinguser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("GetByUsername", mock.Anything, "existinguser").
			Return(conflict, nil)

		_, err := svc.UpdateByID(context.Background(), input)

		require.ErrorIs(t, err, domain.ErrUserNameConflict)
		profileRepo.AssertExpectations(t)
	})

	t.Run("username not conflict when same profile", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "testuser",
		}
		existing := domain.Profile{
			ID:       profileID,
			Username: "testuser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("GetByUsername", mock.Anything, "testuser").
			Return(existing, nil)
		profileRepo.On("UpdateByID", mock.Anything, input).
			Return(*input, nil)

		updated, err := svc.UpdateByID(context.Background(), input)

		require.NoError(t, err)
		assert.Equal(t, input.ID, updated.ID)
		assert.Equal(t, input.Username, updated.Username)

		profileRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "updateduser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrNotFound)
		_, err := svc.UpdateByID(context.Background(), input)
		require.ErrorIs(t, err, domain.ErrNotFound)
		profileRepo.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "updateduser",
		}
		existing := domain.Profile{
			ID:       profileID,
			Username: "testuser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("GetByUsername", mock.Anything, "updateduser").
			Return(domain.Profile{}, domain.ErrNotFound)
		profileRepo.On("UpdateByID", mock.Anything, input).
			Return(domain.Profile{}, domain.ErrInternalServerError)

		_, err := svc.UpdateByID(context.Background(), input)

		require.ErrorIs(t, err, domain.ErrInternalServerError)
		profileRepo.AssertExpectations(t)
	})

	t.Run("internal error on get existing by ID", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "updateduser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrInternalServerError)
		_, err := svc.UpdateByID(context.Background(), input)
		require.ErrorIs(t, err, domain.ErrInternalServerError)
		profileRepo.AssertExpectations(t)
	})

	t.Run("internal error on get by username", func(t *testing.T) {
		svc, profileRepo := setupService(t)

		profileID := uuid.New()
		input := &domain.Profile{
			ID:       profileID,
			Username: "updateduser",
		}
		existing := domain.Profile{
			ID:       profileID,
			Username: "testuser",
		}
		profileRepo.On("GetByID", mock.Anything, profileID).
			Return(existing, nil)
		profileRepo.On("GetByUsername", mock.Anything, "updateduser").
			Return(domain.Profile{}, domain.ErrInternalServerError)
		_, err := svc.UpdateByID(context.Background(), input)
		require.ErrorIs(t, err, domain.ErrInternalServerError)
		profileRepo.AssertExpectations(t)
	})
}

func setupService(t *testing.T) (*profile.Service, *mocks.MockProfileRepository) {
	t.Helper()
	profileRepo := new(mocks.MockProfileRepository)
	svc := profile.NewService(profileRepo, nil)
	return svc, profileRepo
}
