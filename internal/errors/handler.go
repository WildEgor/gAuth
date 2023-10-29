package error_handler

import (
	"errors"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

var ErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	if e != nil {
		log.Error(e.Error())

		c.Status(code).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": e.Error(),
			},
		})
	} else {
		c.Status(code).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "UNKNOWN_ERROR",
			},
		})
	}

	return nil
}
