package routes

import (
	"btwarch/config"
	"btwarch/handlers"
	"btwarch/middleware"
	"btwarch/services"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRouter(app *fiber.App) {
	config := config.LoadConfig()
	authHandler := handlers.NewAuthHandler(config)
	authService := services.NewAuthService(
		config.JWTSecret,
		config.CookieDomain,
		config.CookieSecure,
		config.CookieSameSite,
	)

	authGroup := app.Group("/auth")

	authGroup.Get("/github", authHandler.InitiateGitHubAuth)
	authGroup.Get("/github/callback", authHandler.GitHubCallback)
	authGroup.Post("/logout", authHandler.Logout)
	authGroup.Get("/check", middleware.AuthMiddleware(authService), authHandler.CheckAuth)
}
