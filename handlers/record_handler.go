package handlers

import (
	"btwarch/repositories"

	"github.com/gofiber/fiber/v2"
)

type RecordHandler struct {
	recordRepo *repositories.RecordRepository
}

func NewRecordHandler(recordRepo *repositories.RecordRepository) *RecordHandler {
	return &RecordHandler{recordRepo: recordRepo}
}

func (h *RecordHandler) CreateRecord(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *RecordHandler) GetRecords(c *fiber.Ctx) error {
	// TODO
	return c.JSON(fiber.Map{
		"message": "Records fetched successfully",
	})
}

func (h *RecordHandler) GetRecord(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *RecordHandler) UpdateRecord(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *RecordHandler) DeleteRecord(c *fiber.Ctx) error {
	// TODO
	return nil
}
