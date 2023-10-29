package otp_generate_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPGenHandler struct{}

func NewOTPGenHandler() *OTPGenHandler {
	return &OTPGenHandler{}
}

func (h *OTPGenHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl otp generate logic here

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"identity_type": dto.Identity,
			"code":          "", // for debug
		},
	})

	return nil
}

func (h *OTPGenHandler) parseAndValidate(c *fiber.Ctx) (*authDtos.OTPGenerateRequestDto, error) {
	// Create a new user auth struct.
	dto := &authDtos.OTPGenerateRequestDto{}

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
