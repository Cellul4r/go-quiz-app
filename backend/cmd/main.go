package main

import (
	"log"

	"github.com/Cellul4r/go-quiz-app/api/routes"
	config "github.com/Cellul4r/go-quiz-app/pkg/config"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Load environment variables from .env file
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	app := fiber.New()
	routes.Register(app)

	// Example route to test the server
	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.SendString("Welcome it works! Hello Worlddddd!")
	})

	port := config.Port
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
