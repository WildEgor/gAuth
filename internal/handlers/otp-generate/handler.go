package otp_generate_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
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

	resp := core_dtos.InitResponse()

	us, err := h.ur.FindByPhone(dto.Phone)
	if err != nil {
		resp.SetStatus(c, fiber.StatusUnauthorized)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR: authority", // TODO: make better
		})
		resp.FormResponse()

		return nil
	}

	if !us.IsResendAvailable() {
		resp.SetStatus(c, fiber.StatusTooManyRequests)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR: resend not available", // TODO: make better
		})
		resp.FormResponse()
		return nil
	}

	code, err := h.otps.GenerateAndSMSSend(us.Phone)
	if err != nil {
		resp.SetStatus(c, fiber.StatusTooManyRequests)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR: sms send", // TODO: make better
		})
		resp.FormResponse()
		return nil
	}

	us.UpdateOTP(us.Phone, code)

	err = h.ur.UpdateOTP(us.Id, us.OTP)
	if err != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR: unknown", // TODO: make better
		})
		resp.FormResponse()
		return nil
	}

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(fiber.Map{
		"identity_type": dto.Phone,
		"code":          us.OTP.Code, // for debug
	})
	resp.FormResponse()
	return nil
}
