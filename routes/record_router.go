package routes

import (
	"btwarch/config"
	"btwarch/handlers"
	"btwarch/middleware"
	"btwarch/repositories"
	"btwarch/services"

	"github.com/gofiber/fiber/v2"
)

func InitRecordRouter(app *fiber.App) {
	config := config.LoadConfig()
	recordHandler := handlers.NewRecordHandler(repositories.NewRecordRepository())
	authService := services.NewAuthService(
		config.JWTSecret,
		config.CookieDomain,
		config.CookieSecure,
		config.CookieSameSite,
	)

	recordGroup := app.Group("/records")

	recordGroup.Use(middleware.AuthMiddleware(authService))

	recordGroup.Post("/", recordHandler.CreateRecord)
	recordGroup.Get("/", recordHandler.GetRecords)
	recordGroup.Get("/:id", recordHandler.GetRecord)
	recordGroup.Put("/:id", recordHandler.UpdateRecord)
	recordGroup.Delete("/:id", recordHandler.DeleteRecord)
}
