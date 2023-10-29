package otp_generate_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	userDtos "github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type ChangePasswordHandler struct{}

func NewChangePasswordHandler() *ChangePasswordHandler {
	return &ChangePasswordHandler{}
}

func (h *ChangePasswordHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl change password logic here

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"identity_type": dto.NewPassword,
			"code":          "", // for debug
		},
	})

	return nil
}

func (h *ChangePasswordHandler) parseAndValidate(c *fiber.Ctx) (*userDtos.ChangePasswordRequestDto, error) {
	// Create a new user auth struct.
	dto := &userDtos.ChangePasswordRequestDto{}

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
