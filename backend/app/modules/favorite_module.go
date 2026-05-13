package modules

import (
	"log/slog"

	"github.com/Cellul4r/go-quiz-app/backend/favorite"
	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterFavoriteModule(db *gorm.DB, appLogger *slog.Logger,
	api fiber.Router, config config.Config) {
	favoriteRepo := postgresRepo.NewFavoriteRepository(db, appLogger.With("layer", "repository", "component", "favorite"))

	favoriteService := favorite.NewService(favoriteRepo, appLogger.With("layer", "service", "component", "favorite"))
	rest.NewFavoriteHandler(api, &config, favoriteService)
}
