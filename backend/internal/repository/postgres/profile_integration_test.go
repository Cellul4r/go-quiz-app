// go:build integration

package postgres_test

import (
	"context"
	"testing"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByID_Integration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, repo := setupProfileRepository(t)

		seedProfile := domain.Profile{
			ID:        uuid.New(),
			Username:  "testuer",
			FullName:  "Test User",
			AvatarURL: "http://example.com/avatar.png",
		}

		db.Create(&seedProfile)
		t.Cleanup(func() { db.Unscoped().Delete(&seedProfile) })

		profile, err := repo.GetByID(context.Background(), seedProfile.ID)

		require.NoError(t, err)
		assert.Equal(t, seedProfile.ID, profile.ID)
		assert.Equal(t, seedProfile.Username, profile.Username)
		assert.Equal(t, seedProfile.FullName, profile.FullName)
		assert.Equal(t, seedProfile.AvatarURL, profile.AvatarURL)
	})

	t.Run("not found", func(t *testing.T) {
		_, repo := setupProfileRepository(t)

		profileID := uuid.New()
		_, err := repo.GetByID(context.Background(), profileID)

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestGetByUsername_Integration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, repo := setupProfileRepository(t)

		seedProfile := domain.Profile{
			ID:        uuid.New(),
			Username:  "testuer",
			FullName:  "Test User",
			AvatarURL: "http://example.com/avatar.png",
		}
		db.Create(&seedProfile)
		t.Cleanup(func() { db.Unscoped().Delete(&seedProfile) })

		profile, err := repo.GetByUsername(context.Background(), seedProfile.Username)

		require.NoError(t, err)
		assert.Equal(t, seedProfile.ID, profile.ID)
		assert.Equal(t, seedProfile.Username, profile.Username)
		assert.Equal(t, seedProfile.FullName, profile.FullName)
		assert.Equal(t, seedProfile.AvatarURL, profile.AvatarURL)
	})

	t.Run("not found", func(t *testing.T) {
		_, repo := setupProfileRepository(t)

		username := "nonexistentuser"
		_, err := repo.GetByUsername(context.Background(), username)

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUpdateByID_Integration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, repo := setupProfileRepository(t)

		seedProfile := domain.Profile{
			ID:        uuid.New(),
			Username:  "testuer",
			FullName:  "Test User",
			AvatarURL: "http://example.com/avatar.png",
		}
		db.Create(&seedProfile)
		t.Cleanup(func() { db.Unscoped().Delete(&seedProfile) })

		updateData := domain.Profile{
			ID:        seedProfile.ID,
			Username:  "updateduser",
			FullName:  "Updated User",
			AvatarURL: "http://example.com/updated_avatar.png",
		}

		updatedProfile, err := repo.UpdateByID(context.Background(), &updateData)

		require.NoError(t, err)
		assert.Equal(t, updateData.ID, updatedProfile.ID)
		assert.Equal(t, updateData.Username, updatedProfile.Username)
		assert.Equal(t, updateData.FullName, updatedProfile.FullName)
		assert.Equal(t, updateData.AvatarURL, updatedProfile.AvatarURL)
	})

	t.Run("not found", func(t *testing.T) {
		_, repo := setupProfileRepository(t)

		updateData := domain.Profile{
			ID:        uuid.New(),
			Username:  "updateduser",
			FullName:  "Updated User",
			AvatarURL: "http://example.com/updated_avatar.png",
		}

		_, err := repo.UpdateByID(context.Background(), &updateData)

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("duplicate username conflict", func(t *testing.T) {
		db, repo := setupProfileRepository(t)

		profile1 := domain.Profile{
			ID:       uuid.New(),
			Username: "duplicateuser1",
			FullName: "Duplicate User 1",
		}
		profile2 := domain.Profile{
			ID:       uuid.New(),
			Username: "duplicateuser2",
			FullName: "Duplicate User 2",
		}
		db.Create(&profile1)
		db.Create(&profile2)
		t.Cleanup(func() {
			db.Unscoped().Delete(&profile1)
			db.Unscoped().Delete(&profile2)
		})

		updateData := domain.Profile{
			ID:        profile1.ID,
			Username:  profile2.Username, // change to profile2's username to cause conflict
			FullName:  "Updated User",
			AvatarURL: "http://example.com/updated_avatar.png",
		}

		_, err := repo.UpdateByID(context.Background(), &updateData)

		assert.ErrorIs(t, err, domain.ErrInternalServerError)
	})
}
