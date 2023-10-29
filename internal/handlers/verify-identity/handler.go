package otp_generate_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type VerifyIdentityHandler struct{}

func NewVerifyIdentityHandler() *VerifyIdentityHandler {
	return &VerifyIdentityHandler{}
}

func (h *VerifyIdentityHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl verification logic here

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       dto.Identity,
			"access_token":  "",
			"refresh_token": "",
		},
	})

	return nil
}

func (h *VerifyIdentityHandler) parseAndValidate(c *fiber.Ctx) (*authDtos.VerifyIdentityRequestDto, error) {
	// Create a new user auth struct.
	dto := &authDtos.VerifyIdentityRequestDto{}

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
