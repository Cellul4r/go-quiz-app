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

const baseQuizURL = "/api/v1/quizzes"

func setupQuizHandler(svc rest.QuizService) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: domain.NewStructValidator(),
	})

	handler := &rest.QuizHandler{
		Service: svc,
	}

	api := app.Group(baseQuizURL)

	// Public
	api.Get("", handler.GetAllPublic)

	// Protected
	protected := api.Group("/me", func(c fiber.Ctx) error {
		c.Locals("user", &domain.User{ID: uuid.New()})
		return c.Next()
	})
	protected.Get("", handler.GetMyQuizzes)
	protected.Get("/:id", handler.GetMyQuiz)
	protected.Post("", handler.CreateMyQuiz)
	protected.Patch("/:id", handler.UpdateMyQuiz)
	protected.Delete("/:id", handler.DeleteMyQuiz)

	return app
}

func TestGetMyQuiz(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)
		quizID := uuid.New()
		svc.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(domain.Quiz{
				ID:               quizID,
				Title:            "Test Quiz",
				Description:      "This is a test quiz",
				VisibilityStatus: domain.VisibilityPublic,
			}, nil)

		req := httptest.NewRequest(http.MethodGet, baseQuizURL+"/me/"+quizID.String(), nil)

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]any
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, quizID.String(), body["id"])
		assert.Equal(t, "Test Quiz", body["title"])
		assert.Equal(t, "This is a test quiz", body["description"])
		assert.Equal(t, string(domain.VisibilityPublic), body["visibility_status"])

		svc.AssertExpectations(t)
	})

	t.Run("invalid quiz ID", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		req := httptest.NewRequest(http.MethodGet, baseQuizURL+"/me/invalid-uuid", nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "GetByID")
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(domain.Quiz{}, domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, baseQuizURL+"/me/"+quizID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(domain.Quiz{}, domain.ErrInternalServerError)

		req := httptest.NewRequest(http.MethodGet, baseQuizURL+"/me/"+quizID.String(), nil)
		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})
}

func TestGetAllPublic(t *testing.T) {

	quizzes := []domain.Quiz{
		{
			ID:               uuid.New(),
			AuthorID:         uuid.New(),
			Title:            "Public Quiz 1",
			Description:      "This is the first public quiz",
			VisibilityStatus: domain.VisibilityPublic,
		},
		{
			ID:               uuid.New(),
			AuthorID:         uuid.New(),
			Title:            "Public Quiz 2",
			Description:      "This is the second public quiz",
			VisibilityStatus: domain.VisibilityPublic,
		},
	}

	testcases := []struct {
		name               string
		expectedQuizzes    []domain.Quiz
		mockReturnError    error
		expectedStatusCode int
	}{
		{
			name:               "success",
			expectedQuizzes:    quizzes,
			mockReturnError:    nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "internal error",
			expectedQuizzes:    nil,
			mockReturnError:    domain.ErrInternalServerError,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockQuizService)
			app := setupQuizHandler(svc)

			svc.On("GetAll", mock.Anything, (*uuid.UUID)(nil), true).
				Return(tc.expectedQuizzes, tc.mockReturnError)

			req := httptest.NewRequest(http.MethodGet, baseQuizURL, nil)
			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.mockReturnError == nil {
				var body []map[string]any
				err = json.NewDecoder(resp.Body).Decode(&body)

				require.NoError(t, err)
				require.Len(t, body, len(tc.expectedQuizzes))

				for i, quiz := range tc.expectedQuizzes {
					assert.Equal(t, quiz.ID.String(), body[i]["id"])
					assert.Equal(t, quiz.AuthorID.String(), body[i]["author_id"])
					assert.Equal(t, quiz.Title, body[i]["title"])
					assert.Equal(t, quiz.Description, body[i]["description"])
					assert.Equal(t, string(quiz.VisibilityStatus), body[i]["visibility_status"])
				}
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestGetMyQuizzes(t *testing.T) {

	quizzes := []domain.Quiz{
		{
			ID:               uuid.New(),
			AuthorID:         uuid.New(),
			Title:            "My Quiz 1",
			Description:      "This is the first quiz",
			VisibilityStatus: domain.VisibilityPrivate,
		},
		{
			ID:               uuid.New(),
			AuthorID:         uuid.New(),
			Title:            "My Quiz 2",
			Description:      "This is the second quiz",
			VisibilityStatus: domain.VisibilityPublic,
		},
	}
	testcases := []struct {
		name               string
		expectedQuizzes    []domain.Quiz
		onlyPublic         string
		mockReturnError    error
		expectedStatusCode int
		mockServiceCalled  bool
	}{
		{
			name:               "success",
			expectedQuizzes:    quizzes,
			onlyPublic:         "false",
			mockReturnError:    nil,
			expectedStatusCode: http.StatusOK,
			mockServiceCalled:  true,
		},
		{
			name:               "internal error",
			expectedQuizzes:    nil,
			onlyPublic:         "false",
			mockReturnError:    domain.ErrInternalServerError,
			expectedStatusCode: http.StatusInternalServerError,
			mockServiceCalled:  true,
		},
		{
			name:               "bad request - invalid query param (only_public)", // This case is not applicable for GetMyQuizzes, but we can still test it to ensure it returns bad request
			expectedQuizzes:    nil,
			onlyPublic:         "invalid",
			mockReturnError:    nil, // Service should not be called in this case
			expectedStatusCode: http.StatusBadRequest,
			mockServiceCalled:  false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockQuizService)
			app := setupQuizHandler(svc)

			if tc.mockServiceCalled {
				svc.On("GetAll", mock.Anything, mock.AnythingOfType("*uuid.UUID"), false).
					Return(tc.expectedQuizzes, tc.mockReturnError)
			}

			req := httptest.NewRequest(http.MethodGet, baseQuizURL+"/me?only_public="+tc.onlyPublic, nil)

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.mockServiceCalled && tc.mockReturnError == nil {
				var body []map[string]any
				err = json.NewDecoder(resp.Body).Decode(&body)

				require.NoError(t, err)
				require.Len(t, body, len(tc.expectedQuizzes))

				for i, quiz := range tc.expectedQuizzes {
					assert.Equal(t, quiz.ID.String(), body[i]["id"])
					assert.Equal(t, quiz.AuthorID.String(), body[i]["author_id"])
					assert.Equal(t, quiz.Title, body[i]["title"])
					assert.Equal(t, quiz.Description, body[i]["description"])
					assert.Equal(t, string(quiz.VisibilityStatus), body[i]["visibility_status"])
				}
			}

			svc.AssertExpectations(t)
		})
	}
}

func TestCreateMyQuiz(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("Create", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*domain.Quiz")).
			Return(domain.Quiz{
				ID:               quizID,
				Title:            "New Quiz",
				Description:      "This is a new quiz",
				VisibilityStatus: domain.VisibilityPrivate,
			}, nil)

		reqBody := `{
			"title": "New Quiz",
			"description": "This is a new quiz",
			"visibility_status": "private"
		}`
		req := httptest.NewRequest(http.MethodPost, baseQuizURL+"/me", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var body map[string]any
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, quizID.String(), body["id"])
		assert.Equal(t, "New Quiz", body["title"])
		assert.Equal(t, "This is a new quiz", body["description"])
		assert.Equal(t, string(domain.VisibilityPrivate), body["visibility_status"])
		assert.Equal(t, "author_id", body["author_id"])

		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		svc.On("Create", mock.Anything, mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*domain.Quiz")).
			Return(domain.Quiz{}, domain.ErrInternalServerError)

		reqBody := `{
			"title": "New Quiz",
			"description": "This is a new quiz",
			"visibility_status": "private"
		}`
		req := httptest.NewRequest(http.MethodPost, baseQuizURL+"/me", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	validationTestcases := []struct {
		name               string
		body               map[string]any
		expectedStatusCode int
	}{
		{
			name: "missing title",
			body: map[string]any{
				"description":       "This is a new quiz",
				"visibility_status": "private",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "title too long",
			body: map[string]any{
				"title":             strings.Repeat("a", 101),
				"description":       "This is a new quiz",
				"visibility_status": "private",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "description too long",
			body: map[string]any{
				"title":             "New Quiz",
				"description":       strings.Repeat("a", 201),
				"visibility_status": "private",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "missing visibility status",
			body: map[string]any{
				"title":       "New Quiz",
				"description": "This is a new quiz",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid visibility status",
			body: map[string]any{
				"title":             "New Quiz",
				"description":       "This is a new quiz",
				"visibility_status": "invalid",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range validationTestcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockQuizService)
			app := setupQuizHandler(svc)

			bodyBytes, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, baseQuizURL+"/me", strings.NewReader(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			svc.AssertNotCalled(t, "Create")
		})
	}
}

func TestUpdateMyQuiz(t *testing.T) {
	normalCases := []struct {
		name               string
		mockReturnError    error
		expectedStatusCode int
	}{
		{
			name:               "success",
			mockReturnError:    nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "internal error",
			mockReturnError:    domain.ErrInternalServerError,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "not found",
			mockReturnError:    domain.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, tc := range normalCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockQuizService)
			app := setupQuizHandler(svc)

			quizID := uuid.New()
			svc.On("UpdateByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*domain.Quiz")).
				Return(domain.Quiz{
					ID:               quizID,
					Title:            "Updated Quiz",
					Description:      "This is an updated quiz",
					VisibilityStatus: domain.VisibilityPublic,
				}, tc.mockReturnError)

			reqBody := `{
				"title": "Updated Quiz",
				"description": "This is an updated quiz",
				"visibility_status": "public"
			}`
			req := httptest.NewRequest(http.MethodPatch, baseQuizURL+"/me/"+quizID.String(), strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.mockReturnError == nil {
				var body map[string]any
				err = json.NewDecoder(resp.Body).Decode(&body)
				require.NoError(t, err)

				assert.Equal(t, quizID.String(), body["id"])
				assert.Equal(t, "Updated Quiz", body["title"])
				assert.Equal(t, "This is an updated quiz", body["description"])
				assert.Equal(t, string(domain.VisibilityPublic), body["visibility_status"])
			}

			svc.AssertExpectations(t)

		})
	}

	// Validation error case
	t.Run("invalid quiz ID param", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		reqBody := `{
			"title": "Updated Quiz",
			"description": "This is an updated quiz",
			"visibility_status": "public"
		}`
		req := httptest.NewRequest(http.MethodPatch, baseQuizURL+"/me/invalid-uuid", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})

	validationTestcases := []struct {
		name               string
		body               map[string]any
		expectedStatusCode int
	}{
		{
			name: "title too long",
			body: map[string]any{
				"title":             strings.Repeat("a", 101),
				"description":       "This is an updated quiz",
				"visibility_status": "public",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "description too long",
			body: map[string]any{
				"title":             "Updated Quiz",
				"description":       strings.Repeat("a", 201),
				"visibility_status": "public",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "invalid visibility status",
			body: map[string]any{
				"title":             "Updated Quiz",
				"description":       "This is an updated quiz",
				"visibility_status": "invalid",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range validationTestcases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mocks.MockQuizService)
			app := setupQuizHandler(svc)

			quizID := uuid.New()

			bodyBytes, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, baseQuizURL+"/me/"+quizID.String(), strings.NewReader(string(bodyBytes)))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			svc.AssertNotCalled(t, "UpdateByID")
		})
	}
}

func TestDeleteMyQuiz(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("DeleteByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, baseQuizURL+"/me/"+quizID.String(), nil)

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("DeleteByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(domain.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, baseQuizURL+"/me/"+quizID.String(), nil)

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		quizID := uuid.New()
		svc.On("DeleteByID", mock.Anything, mock.AnythingOfType("uuid.UUID"), quizID).
			Return(domain.ErrInternalServerError)

		req := httptest.NewRequest(http.MethodDelete, baseQuizURL+"/me/"+quizID.String(), nil)

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("invalid quiz ID param", func(t *testing.T) {
		svc := new(mocks.MockQuizService)
		app := setupQuizHandler(svc)

		req := httptest.NewRequest(http.MethodDelete, baseQuizURL+"/me/invalid-uuid", nil)

		resp, err := app.Test(req)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "DeleteByID")
	})
}
