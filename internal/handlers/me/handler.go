package me_handler

import (
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

	authUser := middlewares.ExtractUser(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": user.MeDto{
			ID:     authUser.Id.Hex(),
			Mobile: authUser.Phone,
			Email:  authUser.Email,
		},
	})
}
