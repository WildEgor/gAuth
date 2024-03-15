package otp_generate_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	userDtos "github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type ChangePasswordHandler struct{}

func NewChangePasswordHandler() *ChangePasswordHandler {
	return &ChangePasswordHandler{}
}

func (h *ChangePasswordHandler) Handle(c *fiber.Ctx) error {
	dto := &userDtos.ChangePasswordRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	// TODO: impl change password logic here

	resp := core_dtos.InitResponse()

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(fiber.Map{
		"identity_type": dto.NewPassword,
		"code":          "", // for debug
	})
	resp.FormResponse()

	return nil
}
