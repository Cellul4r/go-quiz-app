package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/mocks"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var requesterID = uuid.New()

func setupFavoriteHandler(svc rest.FavoriteService) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: domain.NewStructValidator(),
	})

	handler := &rest.FavoriteHandler{
		Service: svc,
	}

	api := app.Group("/api/v1")

	protected := api.Group("", func(c fiber.Ctx) error {
		c.Locals("user", &domain.User{ID: requesterID})
		return c.Next()
	})
	protected.Get("/favorites", handler.GetFavorites)
	protected.Post("/favorites", handler.AddFavorite)
	protected.Delete("/favorites", handler.RemoveFavorite)
	return app
}

func TestGetFavorites(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)
		profileID := requesterID

		quiz1 := uuid.New()
		quiz2 := uuid.New()
		authorID := uuid.New()

		svc.On("GetFavoritesByProfileID", mock.Anything, requesterID, profileID).
			Return([]domain.Quiz{
				{
					ID:               quiz1,
					AuthorID:         authorID,
					Title:            "Quiz 1",
					Description:      "Description 1",
					VisibilityStatus: domain.VisibilityPublic,
				},
				{
					ID:               quiz2,
					AuthorID:         authorID,
					Title:            "Quiz 2",
					Description:      "Description 2",
					VisibilityStatus: domain.VisibilityPrivate,
				},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/favorites", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		favorites, ok := body["favorites"].([]any)
		require.True(t, ok)
		require.Len(t, favorites, 2)

		first := favorites[0].(map[string]any)
		second := favorites[1].(map[string]any)

		assert.Equal(t, quiz1.String(), first["id"])
		assert.Equal(t, authorID.String(), first["author_id"])
		assert.Equal(t, "Quiz 1", first["title"])
		assert.Equal(t, "Description 1", first["description"])
		assert.Equal(t, string(domain.VisibilityPublic), first["visibility_status"])

		assert.Equal(t, quiz2.String(), second["id"])
		assert.Equal(t, authorID.String(), second["author_id"])
		assert.Equal(t, "Quiz 2", second["title"])
		assert.Equal(t, "Description 2", second["description"])
		assert.Equal(t, string(domain.VisibilityPrivate), second["visibility_status"])
		svc.AssertExpectations(t)

	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)
		profileID := requesterID

		svc.On("GetFavoritesByProfileID", mock.Anything, requesterID, profileID).
			Return([]domain.Quiz{}, domain.ErrInternalServerError)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/favorites", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)

		profileID := uuid.New()
		quizID := uuid.New()

		svc.On("AddFavorite", mock.Anything, requesterID, profileID, quizID).
			Return(nil)

		reqBody := `{
			"profile_id": "` + profileID.String() + `",
			"quiz_id": "` + quizID.String() + `"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/favorites", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)

		profileID := uuid.New()
		quizID := uuid.New()

		svc.On("AddFavorite", mock.Anything, requesterID, profileID, quizID).
			Return(domain.ErrInternalServerError)

		reqBody := `{
			"profile_id": "` + profileID.String() + `",
			"quiz_id": "` + quizID.String() + `"
		}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/favorites", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	validationTestcases := []struct {
		name               string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "missing profile_id",
			body:               `{"quiz_id": "` + uuid.New().String() + `"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing quiz_id",
			body:               `{"profile_id": "` + uuid.New().String() + `"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range validationTestcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockFavoriteService)
			app := setupFavoriteHandler(svc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/favorites", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			svc.AssertNotCalled(t, "AddFavorite")
		})
	}
}

func TestRemoveFavorite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)

		profileID := uuid.New()
		quizID := uuid.New()

		svc.On("RemoveFavorite", mock.Anything, requesterID, profileID, quizID).
			Return(nil)

		reqBody := `{
			"profile_id": "` + profileID.String() + `",
			"quiz_id": "` + quizID.String() + `"
		}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/favorites", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockFavoriteService)
		app := setupFavoriteHandler(svc)

		profileID := uuid.New()
		quizID := uuid.New()

		svc.On("RemoveFavorite", mock.Anything, requesterID, profileID, quizID).
			Return(domain.ErrInternalServerError)

		reqBody := `{
			"profile_id": "` + profileID.String() + `",
			"quiz_id": "` + quizID.String() + `"
		}`
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/favorites", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	validationTestcases := []struct {
		name               string
		body               string
		expectedStatusCode int
	}{
		{
			name:               "missing profile_id",
			body:               `{"quiz_id": "` + uuid.New().String() + `"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "missing quiz_id",
			body:               `{"profile_id": "` + uuid.New().String() + `"}`,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range validationTestcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockFavoriteService)
			app := setupFavoriteHandler(svc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/favorites", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			svc.AssertNotCalled(t, "RemoveFavorite")
		})
	}
}
