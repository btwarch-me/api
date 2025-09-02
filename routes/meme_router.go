package routes

import (
	"btwarch/config"
	"btwarch/handlers"
	"btwarch/middleware"
	"btwarch/repositories"
	"btwarch/services"

	"github.com/gofiber/fiber/v2"
)

func InitMemeRouter(app *fiber.App) {
	config := config.LoadConfig()
	memeHandler := handlers.NewMemeHandler(repositories.NewMemeRepository())
	authService := services.NewAuthService(
		config.JWTSecret,
		config.CookieDomain,
		config.CookieSecure,
		config.CookieSameSite,
	)

	memeGroup := app.Group("/memes")

	memeGroup.Use(middleware.AuthMiddleware(authService))

	memeGroup.Post("/", memeHandler.CreateMeme)
	memeGroup.Get("/", memeHandler.GetMemes)
	memeGroup.Get("/:id", memeHandler.GetMeme)
	memeGroup.Put("/:id", memeHandler.UpdateMeme)
	memeGroup.Delete("/:id", memeHandler.DeleteMeme)
}
