package middleware
import (
	"strings"
	"github.com/gofiber/fiber/v2"
)

func LinuxOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")
		blocked := []string{"Android", "Macintosh", "iPhone", "iPad", "Windows"}

		for _, b := range blocked {
			if strings.Contains(userAgent, b) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Access denied: only Linux users allowed",
				})
			}
		}

		if !strings.Contains(userAgent, "Linux") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: only Linux users allowed",
			})
		}

		return c.Next()
	}
}