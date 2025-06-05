package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lucadev04/Eventtimer-backend/auth"
	"github.com/lucadev04/Eventtimer-backend/models"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// 🌍 CORS aktivieren – wichtig für getrennte Server
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Oder gezielt z. B. "http://localhost:5173"
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// API-Routen (ohne HTML/Render)
	api := app.Group("/api")

	api.Post("/register", func(c *fiber.Ctx) error {
		payload := new(models.User)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if models.UserExists(payload) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
		}

		result := models.CreateUser(payload)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created"})
	})

	api.Post("/login", func(c *fiber.Ctx) error {
		payload := new(models.User)
		if err := c.BodyParser(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if auth.Login(payload) {
			// JWT optional hier
			return c.JSON(fiber.Map{"message": "Login successful"})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	})

	api.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "This is the dashboard"})
	})

	// Server starten
	log.Fatal(app.Listen(":8081"))
}
