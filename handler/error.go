package handler

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct{}

// NotFound returns 404 response
func (handler *ErrorHandler) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("Not Found")
}
