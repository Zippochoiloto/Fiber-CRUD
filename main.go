package main

import (
	"log"

	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/config"
	"github.com/Zippochoiloto/golang-fiber-basic-todo-app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	//dot env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//config DB
	config.ConnectDB()

	setUpRoutes(app)

	err = app.Listen(":8000")

	if err != nil {
		panic(err)
	}
}

func setUpRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the endpoint",
		})
	})

	api := app.Group("/api")
	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the endpoint",
		})
	})

	routes.TodoRoute(api.Group("/todos"))

}
