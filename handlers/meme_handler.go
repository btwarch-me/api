package handlers

import (
	"btwarch/repositories"

	"github.com/gofiber/fiber/v2"
)

type MemeHandler struct {
	memeRepo *repositories.MemeRepository
}

func NewMemeHandler(memeRepo *repositories.MemeRepository) *MemeHandler {
	return &MemeHandler{memeRepo: memeRepo}
}

func (h *MemeHandler) CreateMeme(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *MemeHandler) GetMemes(c *fiber.Ctx) error {
	// TODO
	return c.JSON(fiber.Map{
		"message": "Memes fetched successfully",
	})
}

func (h *MemeHandler) GetMeme(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *MemeHandler) UpdateMeme(c *fiber.Ctx) error {
	// TODO
	return nil
}

func (h *MemeHandler) DeleteMeme(c *fiber.Ctx) error {
	// TODO
	return nil
}
