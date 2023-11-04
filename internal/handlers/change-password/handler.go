package otp_generate_handler

import (
	"github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type ChangePasswordHandler struct{}

func NewChangePasswordHandler() *ChangePasswordHandler {
	return &ChangePasswordHandler{}
}

func (h *ChangePasswordHandler) Handle(c *fiber.Ctx) error {
	dto := &user.ChangePasswordRequestDto{}
	err := validators.ParseAndValidate(c, dto)
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
