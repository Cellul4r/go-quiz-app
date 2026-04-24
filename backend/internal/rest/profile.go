package rest

import (
	"context"
	"net/http"

	"github.com/Cellul4r/go-quiz-app/backend/domain"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/dto"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

var validate = validator.New()

type ProfileService interface {
	GetByID(ctx context.Context, profileID uuid.UUID) (domain.Profile, error)
	UpdateByID(ctx context.Context, profile *domain.Profile) (domain.Profile, error)
}

type ProfileHandler struct {
	Service ProfileService
}

func NewProfileHandler(api fiber.Router, cfg *config.Config, svc ProfileService) {
	handler := &ProfileHandler{
		Service: svc,
	}

	groupProfile := api.Group("/profiles")
	// Public
	groupProfile.Get("/:id", handler.GetByID)

	// Protected
	protected := groupProfile.Group("", middleware.Protected(cfg))
	protected.Put("/me", handler.UpdateMe)
}

// GetByID retrieves a profile by ID
// @Summary Get profile by ID
// @Description Get a user profile by their UUID
// @Tags profiles
// @Accept json
// @Produce json
// @Param id path string true "Profile ID (UUID)"
// @Success 200 {object} dto.ProfileResponse
// @Failure 400 {object} ResponseError "Invalid profile ID"
// @Failure 404 {object} ResponseError "Profile not found"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /profiles/{id} [get]
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

	return c.Status(fiber.StatusOK).JSON(dto.ToProfileResponse(&profile))
}

// UpdateMe updates the authenticated user's profile
// @Summary Update my profile
// @Description Update the current user's profile information (requires authentication)
// @Tags profiles
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ProfileUpdateRequest true "Update profile request"
// @Success 200 {object} dto.ProfileResponse
// @Failure 400 {object} ResponseError "Invalid request or validation error"
// @Failure 401 {object} ResponseError "Unauthorized"
// @Failure 500 {object} ResponseError "Internal server error"
// @Router /profiles/me [put]
func (h *ProfileHandler) UpdateMe(c fiber.Ctx) error {

	ctx := c.Context()

	// Parse the request body into a Profile struct
	var profile dto.ProfileUpdateRequest
	if err := c.Bind().JSON(&profile); err != nil {
		fiberlog.Error("Failed to bind profile update request: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Message: domain.ErrBadParamInput.Error()})
	}

	// Validate the profile Data
	if valid, err := isProfileValid(&profile); !valid {
		fiberlog.Error("Profile validation failed: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(ResponseError{Message: err.Error()})
	}

	// Update the profile
	updatedProfile, err := h.Service.UpdateByID(ctx, profile.ToDomain(getProfileID(c)))
	if err != nil {
		return c.Status(getStatusCode(err)).JSON(ResponseError{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(dto.ToProfileResponse(&updatedProfile))
}

func isProfileValid(p *dto.ProfileUpdateRequest) (bool, error) {
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
