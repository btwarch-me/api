package handlers

import (
	"btwarch/database"
	"btwarch/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RecordHandler struct {
	recordRepo *repositories.RecordRepository
}

func NewRecordHandler(recordRepo *repositories.RecordRepository) *RecordHandler {
	return &RecordHandler{recordRepo: recordRepo}
}

func (h *RecordHandler) CreateRecord(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	var body struct {
		RecordName  string `json:"record_name"`
		RecordType  string `json:"record_type"`
		RecordValue string `json:"record_value"`
		TTL         int    `json:"ttl"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if body.RecordName == "" || body.RecordType == "" || body.RecordValue == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "record_name, record_type, and record_value are required"})
	}

	// Repository handles Cloudflare creation when is_active=true
	record, err := h.recordRepo.CreateRecord(userID, body.RecordName, body.RecordType, body.RecordValue, body.TTL, body.IsActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(record)
}

func (h *RecordHandler) GetRecords(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	records, err := h.recordRepo.GetRecordsByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(records)
}

func (h *RecordHandler) GetRecord(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	idStr := c.Params("id")
	recordID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid record id"})
	}

	record, err := h.recordRepo.GetRecordByID(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if record == nil || record.UserId != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
	}

	return c.JSON(record)
}

func (h *RecordHandler) UpdateRecord(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	idStr := c.Params("id")
	recordID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid record id"})
	}

	existing, err := h.recordRepo.GetRecordByID(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if existing == nil || existing.UserId != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
	}

	var body struct {
		IsActive *bool `json:"is_active"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if body.IsActive == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "is_active is required"})
	}

	// If activating, push to Cloudflare first
	if *body.IsActive && !existing.IsActive {
		cfRecord := database.Record{
			UserId:      existing.UserId,
			RecordName:  existing.RecordName,
			RecordType:  existing.RecordType,
			RecordValue: existing.RecordValue,
			TTL:         existing.TTL,
			IsActive:    true,
		}
		if _, err := h.recordRepo.AddCloudflareRecord(cfRecord); err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": err.Error()})
		}
	}

	if err := h.recordRepo.UpdateRecordStatus(recordID, *body.IsActive); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	updated, err := h.recordRepo.GetRecordByID(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if updated == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to load updated record"})
	}

	return c.JSON(updated)
}

func (h *RecordHandler) CheckAvailability(c *fiber.Ctx) error {
	var body struct {
		RecordName string `json:"record_name"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if body.RecordName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "record_name is required"})
	}

	record, err := h.recordRepo.RecordExists(body.RecordName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	available := !record

	return c.JSON(fiber.Map{
		"available": available,
	})
}

func (h *RecordHandler) DeleteRecord(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	idStr := c.Params("id")
	recordID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid record id"})
	}

	record, err := h.recordRepo.GetRecordByID(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if record == nil || record.UserId != userID {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "record not found"})
	}

	if err := h.recordRepo.DeleteRecord(recordID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
