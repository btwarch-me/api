package services

import (
	"btwarch/config"
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
)

func credentials() (*cloudinary.Cloudinary, context.Context) {
	cfg := config.LoadConfig()

	cld, err := cloudinary.NewFromParams(
		cfg.CloudinaryCloudName,
		cfg.CloudinaryApiKey,
		cfg.CloudinaryApiSecret,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadFile(c *fiber.Ctx) error {
	cld, ctx := credentials()

	fileHeader, err := c.FormFile("FileUpload")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "File not provided",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer file.Close()

	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "storage",
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Upload failed",
		})
	}

	return c.JSON(fiber.Map{
		"url": resp.SecureURL,
	})
}

func ReadFile(c *fiber.Ctx) error {
	cld, ctx := credentials()
	resp, err := cld.Admin.AssetsByAssetFolder(ctx, admin.AssetsByAssetFolderParams{
		AssetFolder: "storage",
		MaxResults:  100,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "files fetch failed",
		})
	}

	urls := []string{}
	for _, file := range resp.Assets {
		urls = append(urls, file.SecureURL)
	}

	return c.JSON(fiber.Map{
		"files": urls,
	})
}

func DeleteFile(c *fiber.Ctx) error {
	cld, ctx := credentials()
	publicId := c.Query("public_id")
	if publicId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing public_id",
		})
	}

	resp, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Delete failed",
		})
	}

	return c.JSON(fiber.Map{"result": resp.Result})
}
