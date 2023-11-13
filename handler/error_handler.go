package handler

import (
	"github.com/gofiber/fiber/v2"
)

// NotFound returns 404 response
func NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("Not Found")
}
