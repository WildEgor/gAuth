package otp_login_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPLoginHandler struct{}

func NewOTPLoginHandler() *OTPLoginHandler {
	return &OTPLoginHandler{}
}

func (h *OTPLoginHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl otp login logic here

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       dto.Phone,
			"access_token":  "",
			"refresh_token": "",
		},
	})

	return nil
}

func (h *OTPLoginHandler) parseAndValidate(c *fiber.Ctx) (*authDtos.OTPLoginRequestDto, error) {
	// Create a new user auth struct.
	dto := &authDtos.OTPLoginRequestDto{}

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
