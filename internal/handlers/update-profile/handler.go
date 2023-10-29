package otp_generate_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	userDtos "github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type UpdateProfileHandler struct{}

func NewUpdateProfileHandler() *UpdateProfileHandler {
	return &UpdateProfileHandler{}
}

func (h *UpdateProfileHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl update profile logic here

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"id":         dto.FirstName,
			"email":      "",
			"phone":      "",
			"first_name": "",
			"last_name":  "",
		},
	})

	return nil
}

func (h *UpdateProfileHandler) parseAndValidate(c *fiber.Ctx) (*userDtos.UpdateProfileRequestDto, error) {
	// Create a new user auth struct.
	dto := &userDtos.UpdateProfileRequestDto{}

	// Checking received data from JSON body.
	if err := c.BodyParser(dto); err != nil {
		// Return status 400 and error message.
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})
	}

	// Create a new validator
	validate := validators.NewValidator()

	// Validate fields.
	if err := validate.Struct(dto); err != nil {
		// Return, if some fields are not valid.
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
				Errors:  validators.ValidatorErrors(err),
			},
		})
	}

	return dto, nil
}
