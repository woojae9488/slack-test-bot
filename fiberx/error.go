package fiberx

import (
	"github.com/gofiber/fiber/v2"
	app "github.com/woojae9488/slack-test-bot"
)

type errBody struct {
	Code    app.ErrCode `json:"code"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	entity := app.GetErrEntity(app.ErrCodeUndefined)

	if e, ok := err.(*app.ErrWithCode); ok {
		entity = app.GetErrEntity(e.Code)
		status = entity.Status
	} else if e, ok := err.(*fiber.Error); ok {
		status = e.Code
	} else {
		// TODO error logging
	}

	body := errBody{
		Code:    entity.Code,
		Message: entity.Message,
		Detail:  err.Error(),
	}
	return c.Status(status).JSON(body)
}
