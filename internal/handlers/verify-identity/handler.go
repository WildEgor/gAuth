package otp_generate_handler

import (
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type VerifyIdentityHandler struct{}

func NewVerifyIdentityHandler() *VerifyIdentityHandler {
	return &VerifyIdentityHandler{}
}

func (h *VerifyIdentityHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.VerifyIdentityRequestDto{}
	err := validators.ParseAndValidate(c, dto)
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
