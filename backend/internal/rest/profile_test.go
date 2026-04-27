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

const baseURL = "/api/v1/profiles"

func setupProfileHandler(svc rest.ProfileService) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: domain.NewStructValidator(),
	})

	handler := &rest.ProfileHandler{
		Service: svc,
	}

	api := app.Group(baseURL)
	api.Get("/:id", handler.GetByID)

	protected := api.Group("", func(c fiber.Ctx) error {
		c.Locals("user", &domain.User{ID: uuid.New()})
		return c.Next()
	})
	protected.Put("/me", handler.UpdateMe)

	return app
}
func TestGetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		profileID := uuid.New()
		svc.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{
				ID:        profileID,
				Username:  "testuser",
				FullName:  "test user",
				AvatarURL: "http://example.com/avatar.png",
			}, nil)

		req := httptest.NewRequest(http.MethodGet, baseURL+"/"+profileID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)

		assert.Equal(t, profileID.String(), body["id"])
		assert.Equal(t, "testuser", body["username"])
		assert.Equal(t, "test user", body["full_name"])
		assert.Equal(t, "http://example.com/avatar.png", body["avatar_url"])

		svc.AssertExpectations(t)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"/invalid-uuid", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "GetByID")
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		profileID := uuid.New()
		svc.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, baseURL+"/"+profileID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		profileID := uuid.New()
		svc.On("GetByID", mock.Anything, profileID).
			Return(domain.Profile{}, domain.ErrInternalServerError)

		req := httptest.NewRequest(http.MethodGet, baseURL+"/"+profileID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})
}

func TestUpdateMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockProfileService)

		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{
				Username:  "updated",
				FullName:  "Updated Name",
				AvatarURL: "http://example.com/update.png",
			}, nil)

		bodySend := `{
			"username": "updated",
			"full_name": "Updated Name",
			"avatar_url": "http://example.com/update.png"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)

		assert.Equal(t, "updated", body["username"])
		assert.Equal(t, "Updated Name", body["full_name"])
		assert.Equal(t, "http://example.com/update.png", body["avatar_url"])
		svc.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(`invalid json`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{}, domain.ErrInternalServerError)

		bodySend := `{
			"username": "updated",
			"full_name": "Updated Name",
			"avatar_url": "http://example.com/update.png"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{}, domain.ErrNotFound)

		bodySend := `{
			"username": "updated",
			"full_name": "Updated Name",
			"avatar_url": "http://example.com/update.png"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("username too long", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{}, domain.ErrBadParamInput)
		bodySend := `{
			"username": "thisusernameiswaytoolongtobevalid",
			"full_name": "Updated Name",
			"avatar_url": "http://example.com/update.png"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})

	t.Run("full name too long", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{}, domain.ErrBadParamInput)
		bodySend := `{
			"username": "thisusernamevalid",
			"full_name": "This is a very long full name that exceeds the allowed limit",
			"avatar_url": "http://example.com/update.png"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})

	t.Run("avatar URL invalid", func(t *testing.T) {
		svc := new(mocks.MockProfileService)
		app := setupProfileHandler(svc)

		svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("*domain.Profile")).
			Return(domain.Profile{}, domain.ErrBadParamInput)
		bodySend := `{
			"username": "thisusernamevalid",
			"full_name": "This allowed limit",
			"avatar_url": "not-a-url"
		}`
		req := httptest.NewRequest(http.MethodPut, baseURL+"/me", strings.NewReader(bodySend))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})
}
