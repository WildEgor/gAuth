package otp_generate_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPGenHandler struct {
	ur   *repositories.UserRepository
	otps *services.OTPService
}

func NewOTPGenHandler(
	ur *repositories.UserRepository,
	otps *services.OTPService,
) *OTPGenHandler {
	return &OTPGenHandler{
		ur,
		otps,
	}
}

func (h *OTPGenHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.OTPGenerateRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	us, err := h.ur.FindByPhone(dto.Phone)
	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: authority", // TODO: make better
			},
		})
	}

	if !us.IsResendAvailable() {
		c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: resend not available", // TODO: make better
			},
		})
	}

	code, err := h.otps.GenerateAndSMSSend(us.Phone)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: sms send", // TODO: make better
			},
		})
	}

	us.UpdateOTP(us.Phone, code)

	err = h.ur.UpdateOTP(us.Id, us.OTP)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: unknown", // TODO: make better
			},
		})
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"identity_type": dto.Phone,
			"code":          us.OTP.Code, // for debug
		},
	})

	return nil
}
