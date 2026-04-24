package main

import (
	"log"
	"os"

	"github.com/Cellul4r/go-quiz-app/backend/internal/config"
	swagger "github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	fiberlog "github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func configureFiberLog(cfg config.Config) {
	fiberlog.SetOutput(os.Stdout)

	if cfg.AppEnv == "development" || cfg.AppDebug {
		fiberlog.SetLevel(fiberlog.LevelDebug)
		return
	}

	fiberlog.SetLevel(fiberlog.LevelInfo)
}

func configureDevelopmentMiddleware(app *fiber.App, cfg config.Config) {
	if cfg.AppEnv != "development" {
		return
	}

	log.Println("Running in development mode")
	app.Use(logger.New())
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1",
		FilePath: "./swagger_doc/swagger.json",
		Path:     "/docs",
	}))
}
