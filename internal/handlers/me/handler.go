package me_handler

import (
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

	c.JSON(fiber.Map{
		"isOk": true,
		"data": nil,
	})
	return nil
}
