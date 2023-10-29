package registration_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct{}

func NewRegistrationHandler() *RegistrationHandler {
	return &RegistrationHandler{}
}

func (h *RegistrationHandler) Handle(c *fiber.Ctx) error {
	dto, err := h.parseAndValidate(c)
	if err != nil {
		return err
	}

	// TODO: impl registration logic here

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"email": dto.Email,
		},
	})

	return nil
}

func (h *RegistrationHandler) parseAndValidate(c *fiber.Ctx) (*authDtos.RegistrationRequestDto, error) {
	// Create a new user auth struct.
	regDto := &authDtos.RegistrationRequestDto{}

	// Checking received data from JSON body.
	if err := c.BodyParser(regDto); err != nil {
		// Return status 400 and error message.
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})
	}

	// Create a new validator for a RegistrationRequestDto.
	validate := validators.NewValidator()

	// Validate fields.
	if err := validate.Struct(regDto); err != nil {
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

	return regDto, nil
}
