package middleware

import (
	"btwarch/config"

	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware(config *config.Config) fiber.Handler {

	allowedOrigins := config.CORSOrigins

	return func(c *fiber.Ctx) error {

		origin := c.Get("Origin")

		if len(allowedOrigins) > 0 {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
		} else {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-HTTP-Method-Override, Accept")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	}
}
