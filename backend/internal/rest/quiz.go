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

func (h *QuizHandler) GetAllPublic(c fiber.Ctx) error {
	ctx := c.Context()

	quizzes, err := h.Service.GetAll(ctx, nil, true)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToQuizResponseList(&quizzes))
}

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
