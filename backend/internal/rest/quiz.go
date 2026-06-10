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

type QuizService interface {
	GetByID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) (domain.Quiz, error)
	GetAll(ctx context.Context, requesterID *uuid.UUID, onlyPublic bool) ([]domain.Quiz, error)
	Create(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error)
	UpdateByID(ctx context.Context, requesterID uuid.UUID, quiz *domain.Quiz) (domain.Quiz, error)
	DeleteByID(ctx context.Context, requesterID uuid.UUID, quizID uuid.UUID) error
}

type QuizHandler struct {
	Service QuizService
}

func NewQuizHandler(api fiber.Router, cfg *config.Config, svc QuizService) {
	handler := &QuizHandler{
		Service: svc,
	}

	groupQuiz := api.Group("/quizzes")
	// Public
	groupQuiz.Get("", handler.GetAllPublic)

	// Protected
	protected := groupQuiz.Group("/me", middleware.Protected(cfg))
	protected.Get("", handler.GetMyQuizzes)
	protected.Get("/:id", handler.GetMyQuiz)
	protected.Post("", handler.CreateMyQuiz)
	protected.Patch("/:id", handler.UpdateMyQuiz)
	protected.Delete("/:id", handler.DeleteMyQuiz)
}

// GetMyQuiz returns a quiz owned by the authenticated user.
// @Summary Get my quiz by ID
// @Description Get a quiz owned by the authenticated user.
// @Tags quizzes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Success 200 {object} dto.QuizResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me/{id} [get]
func (h *QuizHandler) GetMyQuiz(c fiber.Ctx) error {
	quizID, err := uuid.Parse(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Code:    domain.ErrBadParamInput.Code,
			Message: "invalid quiz id",
		})
	}

	ctx := c.Context()
	profileID := getProfileID(c)

	quiz, err := h.Service.GetByID(ctx, profileID, quizID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuizResponse(&quiz))
}

// GetAllPublic returns public quizzes.
// @Summary Get public quizzes
// @Description List quizzes that are visible to everyone.
// @Tags quizzes
// @Produce json
// @Success 200 {array} dto.QuizResponse
// @Failure 500 {object} ResponseError
// @Router /quizzes [get]
func (h *QuizHandler) GetAllPublic(c fiber.Ctx) error {
	ctx := c.Context()

	quizzes, err := h.Service.GetAll(ctx, nil, true)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuizResponseList(&quizzes))
}

// GetMyQuizzes returns quizzes for the authenticated user.
// @Summary Get my quizzes
// @Description List the authenticated user's quizzes. Use only_public=true to return only public quizzes.
// @Tags quizzes
// @Produce json
// @Security BearerAuth
// @Param only_public query bool false "Only public quizzes"
// @Success 200 {array} dto.QuizResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me [get]
func (h *QuizHandler) GetMyQuizzes(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)

	quizQuery := new(dto.QuizQueryFilter)
	err := c.Bind().Query(quizQuery)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Code:    domain.ErrBadParamInput.Code,
			Message: "invalid query parameter for only_public",
		})
	}
	quizzes, err := h.Service.GetAll(ctx, &profileID, quizQuery.OnlyPublic)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuizResponseList(&quizzes))
}

// CreateMyQuiz creates a quiz for the authenticated user.
// @Summary Create my quiz
// @Description Create a new quiz owned by the authenticated user.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.QuizCreateRequest true "Quiz create request"
// @Success 201 {object} dto.QuizResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me [post]
func (h *QuizHandler) CreateMyQuiz(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)

	var quiz dto.QuizCreateRequest
	if err := c.Bind().Body(&quiz); err != nil {
		if fields, ok := mapValidationErrors(err, &quiz); ok {
			return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
				Code:    "validation_failed",
				Message: "validation failed",
				Fields:  fields,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	createdQuiz, err := h.Service.Create(ctx, profileID, quiz.ToDomain())
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToQuizResponse(&createdQuiz))
}

// UpdateMyQuiz updates a quiz owned by the authenticated user.
// @Summary Update my quiz
// @Description Update an existing quiz owned by the authenticated user.
// @Tags quizzes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param request body dto.QuizUpdateRequest true "Quiz update request"
// @Success 200 {object} dto.QuizResponse
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me/{id} [patch]
func (h *QuizHandler) UpdateMyQuiz(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)

	quizID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Code:    domain.ErrBadParamInput.Code,
			Message: "invalid quiz id",
		})
	}

	var quiz dto.QuizUpdateRequest
	if err := c.Bind().Body(&quiz); err != nil {
		if fields, ok := mapValidationErrors(err, &quiz); ok {
			return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
				Code:    "validation_failed",
				Message: "validation failed",
				Fields:  fields,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	updatedQuiz, err := h.Service.UpdateByID(ctx, profileID, quiz.ToDomain(quizID))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuizResponse(&updatedQuiz))
}

// DeleteMyQuiz deletes a quiz owned by the authenticated user.
// @Summary Delete my quiz
// @Description Delete a quiz owned by the authenticated user.
// @Tags quizzes
// @Produce json
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Success 204
// @Failure 400 {object} ResponseError
// @Failure 401 {object} ResponseError
// @Failure 403 {object} ResponseError
// @Failure 404 {object} ResponseError
// @Failure 500 {object} ResponseError
// @Router /quizzes/me/{id} [delete]
func (h *QuizHandler) DeleteMyQuiz(c fiber.Ctx) error {
	ctx := c.Context()
	profileID := getProfileID(c)

	quizID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Code:    domain.ErrBadParamInput.Code,
			Message: "invalid quiz id",
		})
	}

	err = h.Service.DeleteByID(ctx, profileID, quizID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
