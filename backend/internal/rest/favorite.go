package rest

import (
	"context"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/dto"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/middleware"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

type FavoriteService interface {
	GetFavoritesByProfileID(ctx context.Context, requesterID uuid.UUID, profileID uuid.UUID) ([]string, error)
	AddFavorite(ctx context.Context, requesterID, profileID, quizID uuid.UUID) error
	RemoveFavorite(ctx context.Context, requesterID, profileID, quizID uuid.UUID) error
}

type FavoriteHandler struct {
	Service FavoriteService
}

func NewFavoriteHandler(api fiber.Router, cfg *config.Config, svc FavoriteService) {
	handler := &FavoriteHandler{
		Service: svc,
	}

	groupFavorite := api.Group("/favorites")
	protected := groupFavorite.Group("", middleware.Protected(cfg))
	protected.Get("", handler.GetFavorites)
	protected.Post("", handler.AddFavorite)
	protected.Delete("", handler.RemoveFavorite)
}

func (h *FavoriteHandler) GetFavorites(c fiber.Ctx) error {
	ctx := c.Context()
	requesterID := getProfileID(c)
	profileID := requesterID

	favorites, err := h.Service.GetFavoritesByProfileID(ctx, requesterID, profileID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.JSON(fiber.Map{"favorites": favorites})
}

func (h *FavoriteHandler) AddFavorite(c fiber.Ctx) error {
	ctx := c.Context()

	var req dto.FavoriteRequest
	if err := c.Bind().Body(&req); err != nil {
		fiberlog.Errorf("failed to bind request body: %v", err)
		if fields, ok := mapValidationErrors(err, &req); ok {
			return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
				Code:    "validation_failed",
				Message: "validation failed",
				Fields:  fields,
			})
		}
	}

	requesterID := getProfileID(c)
	profileID := req.ProfileID
	quizID := req.QuizID

	err := h.Service.AddFavorite(ctx, requesterID, profileID, quizID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *FavoriteHandler) RemoveFavorite(c fiber.Ctx) error {
	ctx := c.Context()

	var req dto.FavoriteRequest
	if err := c.Bind().Body(&req); err != nil {
		fiberlog.Errorf("failed to bind request body: %v", err)
		if fields, ok := mapValidationErrors(err, &req); ok {
			return c.Status(fiber.StatusBadRequest).JSON(ResponseError{
				Code:    "validation_failed",
				Message: "validation failed",
				Fields:  fields,
			})
		}
	}

	requesterID := getProfileID(c)
	profileID := req.ProfileID
	quizID := req.QuizID

	err := h.Service.RemoveFavorite(ctx, requesterID, profileID, quizID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(toResponseError(err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
