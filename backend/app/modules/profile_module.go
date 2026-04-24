package modules

import (
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/profile"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterProfileModule(db *gorm.DB, appLogger *slog.Logger,
	api fiber.Router, config config.Config) {
	// Prepare repositories
	profileRepo := postgresRepo.NewProfileRepository(db, appLogger.With("layer", "repository", "component", "profile"))

	// BUild services layer
	profileService := profile.NewService(profileRepo, appLogger.With("layer", "service", "component", "profile"))
	rest.NewProfileHandler(api, &config, profileService)
}
