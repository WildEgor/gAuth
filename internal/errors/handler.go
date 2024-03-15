package error_handler

import (
	"errors"
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	log "github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

var ErrorHandler = func(c *fiber.Ctx, err error) error {
	resp := core_dtos.InitResponse()

	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	headers := make(map[string]string, 0)
	headers[fiber.HeaderContentType] = fiber.MIMEApplicationJSON

	resp.SetHeaders(c, headers)

	if e != nil {
		log.Error(e.Error())

		resp.SetStatus(c, code)
		resp.SetData(fiber.Map{
			"message": e.Error(),
		})
	} else {
		resp.SetData(fiber.Map{
			"message": "UNKNOWN_ERROR",
		})
	}

	resp.FormResponse()
	err = resp.JSON(c)
	if err != nil {
		return err
	}

	return nil
}
