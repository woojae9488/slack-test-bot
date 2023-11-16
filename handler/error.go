package handler

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct{}

func NewErrorHandler() ErrorHandler {
	return ErrorHandler{}
}

// NotFound returns 404 response
func (h *ErrorHandler) NotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("Not Found")
}
