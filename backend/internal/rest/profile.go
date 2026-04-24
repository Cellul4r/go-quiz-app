package rest

import (
	"context"
	"net/http"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var validate = validator.New()

type ProfileService interface {
	GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error)
	UpdateByID(ctx context.Context, profile *domain.Profile) error
}

type ProfileHandler struct {
	Service ProfileService
}

func NewProfileHandler(app *fiber.App, cfg *config.Config, svc ProfileService) {
	handler := &ProfileHandler{
		Service: svc,
	}

	api := app.Group("/api/v1/profiles")
	// Public
	api.Get("/:id", handler.GetByID)

	// Protected
	protected := api.Group("", middleware.Protected(cfg))
	protected.Put("/me", handler.UpdateMe)
}

func (h *ProfileHandler) GetByID(c fiber.Ctx) error {
	profileID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(ResponseError{
			Message: "invalid profile id",
		})
	}

	ctx := c.Context()
	profile, err := h.Service.GetByID(ctx, profileID)
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ResponseError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

func (h *ProfileHandler) UpdateMe(c fiber.Ctx) error {

	ctx := c.Context()

	// Parse the request body into a Profile struct
	var profile domain.Profile
	if err := c.Bind().JSON(&profile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Message: domain.ErrBadParamInput.Error()})
	}

	// Validate the profile Data
	if valid, err := isProfileValid(&profile); !valid {
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Message: err.Error()})
	}

	profile.ID = getProfileID(c)

	// Update the profile
	if err := h.Service.UpdateByID(ctx, &profile); err != nil {
		return c.Status(getStatusCode(err)).JSON(ResponseError{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(profile)
}

func isProfileValid(p *domain.Profile) (bool, error) {
	err := validate.Struct(p)
	if err != nil {
		return false, err
	}
	return true, nil
}

// helpers
func getProfileID(c fiber.Ctx) uuid.UUID {
	user := c.Locals("user").(*domain.User)
	id, _ := uuid.Parse(user.ID)
	return id
}
