package main

import (
	"log"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	"github.com/Cellul4r/go-quiz-app/backend/internal/database"
	postgresRepo "github.com/Cellul4r/go-quiz-app/backend/internal/repository/postgres"
	"github.com/Cellul4r/go-quiz-app/backend/internal/rest"
	"github.com/Cellul4r/go-quiz-app/backend/profile"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
)

func main() {
	// Load environment variables from .env file
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	db, dbErr := database.ConnectDatabase(config)
	if dbErr != nil {
		log.Fatal("cannot connect to database: ", dbErr)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	// Defer the closing of the database connection pool
	defer sqlDB.Close()

	app := fiber.New()
	configureFiberLog(config)
	configureDevelopmentMiddleware(app, config)
	appLogger := newAppLogger(config)

	api := app.Group("/api/v1")

	// Prepare repositories
	profileRepo := postgresRepo.NewProfileRepository(db, appLogger.With("layer", "repository", "component", "profile"))

	// BUild services layer
	profileService := profile.NewService(profileRepo, appLogger.With("layer", "service", "component", "profile"))
	rest.NewProfileHandler(api, &config, profileService)
	fiberlog.Info("Server is running on port " + config.Port)
	port := config.Port
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
