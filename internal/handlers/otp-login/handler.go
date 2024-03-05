package otp_login_handler

import (
	"github.com/WildEgor/gAuth/internal/configs"
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type OTPLoginHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewOTPLoginHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *OTPLoginHandler {
	return &OTPLoginHandler{
		ur,
		tr,
		jwt,
		jwtConfig,
	}
}

func (h *OTPLoginHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.OTPLoginRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	// TODO: impl otp login logic here
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

	if !us.VerifyOTP(us.Phone, dto.Code) {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: authority", // TODO: make better
			},
		})
	}

	us.ClearOTP()
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

	// 3. Generate tokens
	at, atErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(us.Id.Hex(), h.jwtConfig.ATDuration)
	if atErr != nil || rtErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: tokens", // TODO: make better
			},
		})

		return nil
	}

	errAT := h.tr.SetAT(at)
	errRT := h.tr.SetRT(rt)
	if errAT != nil || errRT != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR", // TODO: make better
			},
		})

		return nil
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       us.Id.Hex(),
			"access_token":  at.Token,
			"refresh_token": rt.Token,
		},
	})

	return nil
}
