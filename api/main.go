package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alanloffler/shorten-url-fiber-redis/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", routes.GetAll)
	app.Get("/:url", routes.ResolveURL)
	app.Get("/api/rate-limit", routes.GetRateLimit)
	app.Post("/api", routes.ShortenURL)
	app.Delete("/api/:url", routes.DeleteURL)
	app.Delete("/api/rate-limit", routes.ClearRateLimit)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
