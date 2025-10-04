package handlers

import (
	"btwarch/config"
	"btwarch/database"
	"btwarch/repositories"
	"btwarch/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RecordHandler struct {
	recordRepo         *repositories.RecordRepository
	subdomainClaimRepo *repositories.SubdomainClaimRepository
}

func NewRecordHandler(recordRepo *repositories.RecordRepository, subdomainClaimRepo *repositories.SubdomainClaimRepository) *RecordHandler {
	return &RecordHandler{
		recordRepo:         recordRepo,
		subdomainClaimRepo: subdomainClaimRepo,
	}
}

func (h *RecordHandler) ClaimSubdomain(c *fiber.Ctx) error {
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
		SubdomainName string `json:"subdomain_name"`
	}

	if err := c.BodyParser(&body); err != nil || body.SubdomainName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "subdomain_name is required"})
	}

	if err := utils.ValidateSubdomainName(body.SubdomainName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingUserClaim, err := h.subdomainClaimRepo.GetClaimByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if existingUserClaim != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user already has a subdomain claim. Only one subdomain per user is allowed"})
	}

	existingClaim, err := h.subdomainClaimRepo.GetClaimBySubdomain(body.SubdomainName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if existingClaim != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "subdomain already claimed"})
	}

	claim, err := h.subdomainClaimRepo.CreateClaim(userID, body.SubdomainName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "subdomain claimed successfully",
		"claim":       claim,
		"full_domain": utils.GetFullSubdomainName(body.SubdomainName),
	})
}

func (h *RecordHandler) DeleteSubdomain(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claim, err := h.subdomainClaimRepo.GetClaimByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if claim == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no subdomain claim found"})
	}

	if err := h.subdomainClaimRepo.DeleteClaim(claim.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "subdomain claim deleted successfully",
		"claim":   claim,
	})
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

	config := config.LoadConfig()
	if !strings.HasSuffix(body.RecordName, "."+config.ParentDomain) {
		body.RecordName = body.RecordName + "." + config.ParentDomain
	}

	subdomainName := utils.ExtractSubdomainFromRecordName(body.RecordName)
	if subdomainName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid record name format"})
	}

	if body.RecordType != "TXT" {
		claim, err := h.subdomainClaimRepo.GetClaimBySubdomain(subdomainName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if claim == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "subdomain not claimed. Please claim the subdomain first"})
		}

		if claim.UserId != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "subdomain claimed by another user"})
		}
	}

	if err := utils.ValidateRecordName(body.RecordName, body.RecordType, subdomainName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	existingRecord, err := h.recordRepo.GetRecordByNameAndType(body.RecordName, body.RecordType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if body.RecordType == "TXT" {
		body.RecordValue = fmt.Sprintf(`"%s"`, body.RecordValue)
	}

	if existingRecord != nil {
		if err := h.recordRepo.UpdateRecord(existingRecord.ID, body.RecordName, body.RecordType, body.RecordValue, body.TTL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if body.IsActive && existingRecord.CloudflareRecordID != nil {
			cfRecord := database.Record{
				UserId:      userID,
				RecordName:  body.RecordName,
				RecordType:  body.RecordType,
				RecordValue: body.RecordValue,
				TTL:         body.TTL,
				IsActive:    true,
			}
			_, err := h.recordRepo.UpdateOnCloudflare(*existingRecord.CloudflareRecordID, cfRecord)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": utils.ExtractErrorMessage(err),
				})
			}
		}

		updatedRecord, err := h.recordRepo.GetRecordByID(existingRecord.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": utils.ExtractErrorMessage(err),
			})
		}

		return c.Status(fiber.StatusOK).JSON(updatedRecord)
	}

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

	if records == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no records found"})
	}

	return c.JSON(fiber.Map{
		"records": records,
		"message": "records fetched successfully",
	})
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
		RecordName         string `json:"record_name"`
		RecordType         string `json:"record_type"`
		RecordValue        string `json:"record_value"`
		TTL                int    `json:"ttl"`
		IsActive           bool   `json:"is_active"`
		CloudflareRecordID string `json:"cloudflare_record_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if !body.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "is_active is required and must be true	"})
	}

	config := config.LoadConfig()

	if !strings.HasSuffix(body.RecordName, "."+config.ParentDomain) {
		body.RecordName = body.RecordName + "." + config.ParentDomain
	}

	if body.RecordType == "TXT" {
		body.RecordValue = fmt.Sprintf(`"%s"`, body.RecordValue)
	}

	if body.IsActive {
		cfRecord := database.Record{
			UserId:      existing.UserId,
			RecordName:  body.RecordName,
			RecordType:  body.RecordType,
			RecordValue: body.RecordValue,
			TTL:         body.TTL,
			IsActive:    true,
		}

		if err := h.recordRepo.UpdateRecord(recordID, body.RecordName, body.RecordType, body.RecordValue, body.TTL); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if existing.CloudflareRecordID != nil {
			if _, err := h.recordRepo.UpdateOnCloudflare(*existing.CloudflareRecordID, cfRecord); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			newCfID, err := h.recordRepo.CreateCloudflareRecord(cfRecord)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}

			if err := h.recordRepo.UpdateCloudflareIDByNameAndType(body.RecordName, body.RecordType, newCfID.ID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		if err := h.recordRepo.UpdateRecordStatus(recordID, true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		if existing.CloudflareRecordID != nil {
			if err := h.recordRepo.DeleteCloudflareRecord(*existing.CloudflareRecordID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
		}

		if err := h.recordRepo.UpdateRecordStatus(recordID, false); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	updated, err := h.recordRepo.GetRecordByID(recordID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

	config := config.LoadConfig()
	if !strings.HasSuffix(body.RecordName, "."+config.ParentDomain) {
		body.RecordName = body.RecordName + "." + config.ParentDomain
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "record deleted successfully",
		"record":  record,
	})
}

func (h *RecordHandler) GetSubdomainClaim(c *fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	userIDStr, ok := userIDVal.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claim, err := h.subdomainClaimRepo.GetClaimByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if claim == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no subdomain claim found"})
	}

	return c.JSON(claim)
}
