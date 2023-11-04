package otp_generate_handler

import (
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPGenHandler struct{}

func NewOTPGenHandler() *OTPGenHandler {
	return &OTPGenHandler{}
}

func (h *OTPGenHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.OTPGenerateRequestDto{}
	err := validators.ParseAndValidate(c, dto)
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
