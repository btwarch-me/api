package routes

import (
	"btwarch/handlers"
	"btwarch/repositories"

	"github.com/gofiber/fiber/v2"
)

func InitRecordRouter(app *fiber.App) {
	recordHandler := handlers.NewRecordHandler(repositories.NewRecordRepository())

	recordGroup := app.Group("/records")

	recordGroup.Post("/", recordHandler.CreateRecord)
	recordGroup.Get("/", recordHandler.GetRecords)
	recordGroup.Get("/:id", recordHandler.GetRecord)
	recordGroup.Put("/:id", recordHandler.UpdateRecord)
	recordGroup.Delete("/:id", recordHandler.DeleteRecord)
}
