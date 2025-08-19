package routes

import (
	"btwarch/config"
	"btwarch/handlers"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRouter(app *fiber.App) {
	config := config.LoadConfig()
	authHandler := handlers.NewAuthHandler(config)

	authGroup := app.Group("/auth")

	authGroup.Get("/github", authHandler.InitiateGitHubAuth)
	authGroup.Get("/github/callback", authHandler.GitHubCallback)
}
