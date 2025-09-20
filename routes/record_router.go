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
	recordHandler := handlers.NewRecordHandler(
		repositories.NewRecordRepository(),
		repositories.NewSubdomainClaimRepository(),
	)
	authService := services.NewAuthService(
		config.JWTSecret,
		config.CookieDomain,
		config.CookieSecure,
		config.CookieSameSite,
	)

	recordGroup := app.Group("/records")

	recordGroup.Use(middleware.AuthMiddleware(authService))

	recordGroup.Post("/", recordHandler.CreateRecord)
	recordGroup.Post("/claim", recordHandler.ClaimRecord)
	recordGroup.Get("/", recordHandler.GetRecords)
	recordGroup.Get("/claim", recordHandler.GetSubdomainClaim)
	recordGroup.Get("/:id", recordHandler.GetRecord)
	recordGroup.Put("/:id", recordHandler.UpdateRecord)
	recordGroup.Delete("/:id", recordHandler.DeleteRecord)

	recordGroup.Post("/checkavailability", recordHandler.CheckAvailability)

}
