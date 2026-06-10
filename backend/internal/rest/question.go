package rest

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/dto"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type QuestionService interface {
	GetByID(ctx context.Context, requesterID uuid.UUID, questionID string) (*domain.Question, error)
	GetByQuizID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) ([]domain.Question, error)
	Create(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error)
	UpdateByID(ctx context.Context, requesterID uuid.UUID, question *domain.Question) (*domain.Question, error)
	DeleteByID(ctx context.Context, requesterID uuid.UUID, questionID string) error
}

type QuestionHandler struct {
	Service QuestionService
}

func NewQuestionHandler(api fiber.Router, cfg *config.Config, svc QuestionService) {
	handler := &QuestionHandler{Service: svc}

	api.Get("/questions/:id", handler.GetPublicQuestion)
	api.Get("/quizzes/:quiz_id/questions", handler.GetPublicQuestionsByQuiz)

	questionGroup := api.Group("/questions/me", middleware.Protected(cfg))
	questionGroup.Get("/:id", handler.GetMyQuestion)
	questionGroup.Patch("/:id", handler.UpdateMyQuestion)
	questionGroup.Delete("/:id", handler.DeleteMyQuestion)

	quizGroup := api.Group("/quizzes/me/:quiz_id/questions", middleware.Protected(cfg))
	quizGroup.Get("", handler.GetMyQuestionsByQuiz)
	quizGroup.Post("", handler.CreateMyQuestion)
}

// GetPublicQuestion retrieves a public question by ID.
// @Summary Get public question by ID
// @Tags questions
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /questions/{id} [get]
func (h *QuestionHandler) GetPublicQuestion(c fiber.Ctx) error {
	return h.getQuestion(c, uuid.Nil)
}

// GetMyQuestion retrieves a question by ID for the authenticated user.
// @Summary Get my question by ID
// @Tags questions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Question ID"
// @Success 200 {object} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /questions/me/{id} [get]
func (h *QuestionHandler) GetMyQuestion(c fiber.Ctx) error {
	return h.getQuestion(c, getProfileID(c))
}

// GetPublicQuestionsByQuiz retrieves questions for a public quiz.
// @Summary Get public quiz questions
// @Tags questions
// @Produce json
// @Param quiz_id path string true "Quiz ID"
// @Success 200 {array} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/{quiz_id}/questions [get]
func (h *QuestionHandler) GetPublicQuestionsByQuiz(c fiber.Ctx) error {
	return h.getQuestionsByQuiz(c, uuid.Nil)
}

// GetMyQuestionsByQuiz retrieves questions for an authenticated user's quiz.
// @Summary Get my quiz questions
// @Tags questions
// @Produce json
// @Security BearerAuth
// @Param quiz_id path string true "Quiz ID"
// @Success 200 {array} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me/{quiz_id}/questions [get]
func (h *QuestionHandler) GetMyQuestionsByQuiz(c fiber.Ctx) error {
	return h.getQuestionsByQuiz(c, getProfileID(c))
}

// CreateMyQuestion creates a question in an authenticated user's quiz.
// @Summary Create my quiz question
// @Tags questions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param quiz_id path string true "Quiz ID"
// @Param request body dto.QuestionCreateRequest true "Question create request"
// @Success 201 {object} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me/{quiz_id}/questions [post]
func (h *QuestionHandler) CreateMyQuestion(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)

	quizID, err := parseUUIDParam(c, "quiz_id", "invalid quiz id")
	if err != nil {
		return err
	}

	var request dto.QuestionCreateRequest
	if err := c.Bind().Body(&request); err != nil {
		return validationErrorResponse(c, err, &request)
	}

	question, err := request.ToDomain(quizID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: err.Error()})
	}

	created, err := h.Service.Create(ctx, profileID, question)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToQuestionResponse(created))
}

// UpdateMyQuestion updates a question owned by the authenticated user.
// @Summary Update my question
// @Tags questions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Question ID"
// @Param request body dto.QuestionUpdateRequest true "Question update request"
// @Success 200 {object} dto.QuestionResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /questions/me/{id} [patch]
func (h *QuestionHandler) UpdateMyQuestion(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)
	questionID := c.Params("id")

	if _, err := uuid.Parse(questionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: "invalid question id"})
	}

	var request dto.QuestionUpdateRequest
	if err := c.Bind().Body(&request); err != nil {
		return validationErrorResponse(c, err, &request)
	}

	question, err := request.ToDomain(questionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: err.Error()})
	}

	updated, err := h.Service.UpdateByID(ctx, profileID, question)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuestionResponse(updated))
}

// DeleteMyQuestion deletes a question owned by the authenticated user.
// @Summary Delete my question
// @Tags questions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Question ID"
// @Success 204
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /questions/me/{id} [delete]
func (h *QuestionHandler) DeleteMyQuestion(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)
	questionID := c.Params("id")

	if _, err := uuid.Parse(questionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: "invalid question id"})
	}

	if err := h.Service.DeleteByID(ctx, profileID, questionID); err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *QuestionHandler) getQuestion(c fiber.Ctx, requesterID uuid.UUID) error {
	questionID := c.Params("id")
	if _, err := uuid.Parse(questionID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: "invalid question id"})
	}

	question, err := h.Service.GetByID(c.Context(), requesterID, questionID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuestionResponse(question))
}

func (h *QuestionHandler) getQuestionsByQuiz(c fiber.Ctx, requesterID uuid.UUID) error {
	quizID, err := parseUUIDParam(c, "quiz_id", "invalid quiz id")
	if err != nil {
		return err
	}

	questions, err := h.Service.GetByQuizID(c.Context(), requesterID, quizID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuestionResponseList(questions))
}

func parseUUIDParam(c fiber.Ctx, name string, message string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(name))
	if err != nil {
		return uuid.Nil, c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: domain.ErrBadParamInput.Code, Message: message})
	}
	return id, nil
}

func validationErrorResponse(c fiber.Ctx, err error, request any) error {
	if fields, ok := mapValidationErrors(err, request); ok {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Code: "validation_failed", Message: "validation failed", Fields: fields})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
}
