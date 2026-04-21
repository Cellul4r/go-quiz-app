package main

import (
	"log"
	"os"

	"github.com/Cellul4r/go-quiz-app/api/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	app := fiber.New()
	routes.Register(app)

	// Example route to test the server
	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.SendString("Welcome it works! Hello Worlddddd!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
