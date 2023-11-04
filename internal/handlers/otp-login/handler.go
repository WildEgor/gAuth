package otp_login_handler

import (
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPLoginHandler struct{}

func NewOTPLoginHandler() *OTPLoginHandler {
	return &OTPLoginHandler{}
}

func (h *OTPLoginHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.OTPLoginRequestDto{}
	err := validators.ParseAndValidate(c, dto)
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
