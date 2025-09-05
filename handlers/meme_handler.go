package handlers

import (
	"btwarch/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MemeHandler struct {
	memeRepo *repositories.MemeRepository
}

func NewMemeHandler(memeRepo *repositories.MemeRepository) *MemeHandler {
	return &MemeHandler{memeRepo: memeRepo}
}

func (h *MemeHandler) CreateMeme(c *fiber.Ctx) error {
	type MemeRequest struct {
		UserID      string   `json:"user_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	}

	var req MemeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	meme, err := h.memeRepo.CreateMeme(userID, req.Title, req.Description, req.Images)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(meme)
}

func (h *MemeHandler) GetMemes(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	memes, err := h.memeRepo.GetMemesByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(memes)
}

func (h *MemeHandler) GetMeme(c *fiber.Ctx) error {
	id := c.Params("id")
	memeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid meme ID"})
	}

	meme, err := h.memeRepo.GetMemeByID(memeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if meme == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Meme not found"})
	}

	return c.JSON(meme)
}

func (h *MemeHandler) UpdateMeme(c *fiber.Ctx) error {
	id := c.Params("id")
	memeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid meme ID"})
	}

	type MemeRequest struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Images      []string `json:"images"`
	}

	var req MemeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err = h.memeRepo.UpdateMeme(memeID, req.Title, req.Description, req.Images)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Meme updated successfully"})
}

func (h *MemeHandler) DeleteMeme(c *fiber.Ctx) error {
	id := c.Params("id")
	memeID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid meme ID"})
	}

	err = h.memeRepo.DeleteMeme(memeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Meme deleted successfully"})
}
