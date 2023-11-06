package middlewares

import (
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/gofiber/fiber/v2"
)

func ExtractUser(ctx *fiber.Ctx) *models.UsersModel {
	user := ctx.Locals("user").(models.UsersModel)
	return &user
}
