package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

const baseQuestionURL = "/api/v1"

func setupQuestionHandler(svc rest.QuestionService) *fiber.App {
	app := fiber.New(fiber.Config{StructValidator: domain.NewStructValidator()})
	handler := &rest.QuestionHandler{Service: svc}
	api := app.Group(baseQuestionURL)

	api.Get("/questions/:id", handler.GetPublicQuestion)
	api.Get("/quizzes/:quiz_id/questions", handler.GetPublicQuestionsByQuiz)

	protected := func(c fiber.Ctx) error {
		c.Locals("user", &domain.User{ID: testRequesterID})
		return c.Next()
	}
	questionGroup := api.Group("/questions/me", protected)
	questionGroup.Get("/:id", handler.GetMyQuestion)
	questionGroup.Patch("/:id", handler.UpdateMyQuestion)
	questionGroup.Delete("/:id", handler.DeleteMyQuestion)

	quizGroup := api.Group("/quizzes/me/:quiz_id/questions", protected)
	quizGroup.Get("", handler.GetMyQuestionsByQuiz)
	quizGroup.Post("", handler.CreateMyQuestion)

	return app
}

func TestGetPublicQuestion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		questionID := uuid.NewString()
		question := validRESTQuestion(questionID, uuid.New())

		svc.On("GetByID", mock.Anything, uuid.Nil, questionID).Return(question, nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, baseQuestionURL+"/questions/"+questionID, nil))

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, questionID, body["id"])
		assert.Equal(t, string(domain.MultipleChoice), body["question_type"])
		svc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, baseQuestionURL+"/questions/bad-id", nil))

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "GetByID")
	})

	t.Run("forbidden private parent", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		questionID := uuid.NewString()

		svc.On("GetByID", mock.Anything, uuid.Nil, questionID).Return((*domain.Question)(nil), domain.ErrForbidden)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, baseQuestionURL+"/questions/"+questionID, nil))

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		svc.AssertExpectations(t)
	})
}

func TestGetPublicQuestionsByQuiz(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		quizID := uuid.New()
		questions := []domain.Question{*validRESTQuestion(uuid.NewString(), quizID)}

		svc.On("GetByQuizID", mock.Anything, uuid.Nil, quizID).Return(questions, nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, baseQuestionURL+"/quizzes/"+quizID.String()+"/questions", nil))

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		var body []map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		require.Len(t, body, 1)
		assert.Equal(t, questions[0].ID, body[0]["id"])
		svc.AssertExpectations(t)
	})
}

func TestCreateMyQuestion(t *testing.T) {
	t.Run("success multiple choice", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		quizID := uuid.New()
		created := validRESTQuestion(uuid.NewString(), quizID)
		body := `{"question_type":"multiple_choice","content":"What is A?","time_limit_seconds":10,"sort_order":1,"options":{"choices":[{"id":"a","text":"A"},{"id":"b","text":"B"}],"correct_choices":["a"]}}`

		svc.On("Create", mock.Anything, testRequesterID, mock.MatchedBy(func(q *domain.Question) bool {
			_, ok := q.Options.(domain.MultipleChoiceOptions)
			return q.QuizID == quizID.String() && q.Content == "What is A?" && ok
		})).Return(created, nil)

		resp, err := app.Test(jsonRequest(http.MethodPost, baseQuestionURL+"/quizzes/me/"+quizID.String()+"/questions", body))

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		quizID := uuid.New()
		body := `{"question_type":"multiple_choice","content":"missing options"}`

		resp, err := app.Test(jsonRequest(http.MethodPost, baseQuestionURL+"/quizzes/me/"+quizID.String()+"/questions", body))

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "Create")
	})
}

func TestUpdateMyQuestion(t *testing.T) {
	t.Run("success typed answer", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		questionID := uuid.NewString()
		updated := validRESTQuestion(questionID, uuid.New())
		updated.QuestionType = domain.TypedAnswer
		updated.Options = domain.TypedAnswerOptions{CorrectAnswer: "blue"}
		body := `{"question_type":"typed_answer","content":"Updated","options":{"correct_answer":"blue","accepted_answers":["Blue"],"case_sensitive":false}}`

		svc.On("UpdateByID", mock.Anything, testRequesterID, mock.MatchedBy(func(q *domain.Question) bool {
			opts, ok := q.Options.(domain.TypedAnswerOptions)
			return q.ID == questionID && q.QuestionType == domain.TypedAnswer && q.Content == "Updated" && ok && opts.CorrectAnswer == "blue"
		})).Return(updated, nil)

		resp, err := app.Test(jsonRequest(http.MethodPatch, baseQuestionURL+"/questions/me/"+questionID, body))

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("options require question type", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		questionID := uuid.NewString()
		body := `{"options":{"correct_answer":true}}`

		resp, err := app.Test(jsonRequest(http.MethodPatch, baseQuestionURL+"/questions/me/"+questionID, body))

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertNotCalled(t, "UpdateByID")
	})
}

func TestDeleteMyQuestion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := new(mocks.MockQuestionService)
		app := setupQuestionHandler(svc)
		questionID := uuid.NewString()

		svc.On("DeleteByID", mock.Anything, testRequesterID, questionID).Return(nil)

		resp, err := app.Test(httptest.NewRequest(http.MethodDelete, baseQuestionURL+"/questions/me/"+questionID, nil))

		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		svc.AssertExpectations(t)
	})
}

func jsonRequest(method string, target string, body string) *http.Request {
	req := httptest.NewRequest(method, target, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func validRESTQuestion(questionID string, quizID uuid.UUID) *domain.Question {
	return &domain.Question{
		ID:               questionID,
		QuizID:           quizID.String(),
		QuestionType:     domain.MultipleChoice,
		Content:          "question content",
		TimeLimitSeconds: 10,
		SortOrder:        1,
		Options: domain.MultipleChoiceOptions{
			Choices: []domain.Choice{
				{ID: "a", Text: "A"},
				{ID: "b", Text: "B"},
			},
			CorrectChoices: []string{"a"},
		},
	}
}
