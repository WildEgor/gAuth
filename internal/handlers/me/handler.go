package me_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	"github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type MeHandler struct {
	userRepository *repositories.UserRepository
}

func NewMeHandler(
	userRepository *repositories.UserRepository,
) *MeHandler {
	return &MeHandler{
		userRepository,
	}
}

func (hch *MeHandler) Handle(c *fiber.Ctx) error {

	resp := core_dtos.InitResponse()

	authUser := middlewares.ExtractUser(c)

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(user.MeDto{
		ID:     authUser.Id.Hex(),
		Mobile: authUser.Phone,
		Email:  authUser.Email,
	})
	resp.FormResponse()

	return nil
}
