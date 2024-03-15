package otp_generate_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
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

	resp := core_dtos.InitResponse()

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(fiber.Map{
		"id":         dto.FirstName,
		"email":      "",
		"phone":      "",
		"first_name": "",
		"last_name":  "",
	})
	resp.FormResponse()

	// TODO: impl update profile logic here

	return nil
}
