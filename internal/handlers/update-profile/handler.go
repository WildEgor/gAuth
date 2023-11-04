package otp_generate_handler

import (
	userDtos "github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type UpdateProfileHandler struct{}

func NewUpdateProfileHandler() *UpdateProfileHandler {
	return &UpdateProfileHandler{}
}

func (h *UpdateProfileHandler) Handle(c *fiber.Ctx) error {
	dto := &userDtos.UpdateProfileRequestDto{}
	err := validators.ParseAndValidate(c, dto)
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
