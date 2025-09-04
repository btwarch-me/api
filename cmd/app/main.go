package main

import (
	"btwarch/config"
	"btwarch/database"
	// "btwarch/middleware"/
	"btwarch/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting BTWArch API")

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := config.LoadConfig()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.InitTables(); err != nil {
		log.Fatalf("Failed to initialize database tables: %v", err)
	}

	app := fiber.New()

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	app.Use(logger.New())
	// app.Use(middleware.LinuxOnlyMiddleware())

	routes.InitAuthRouter(app)
	routes.InitRecordRouter(app)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
