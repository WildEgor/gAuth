package refresh_handler

import (
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type RefreshHandler struct {
	userRepo *repositories.UserRepository
}

func NewRefreshHandler(
	userRepo *repositories.UserRepository,
) *RefreshHandler {
	return &RefreshHandler{
		userRepo: userRepo,
	}
}

// Handle RefreshHandler method refreshes access token using refresh token.
// @Description Refresh tokens.
// @Summary refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "RefreshToken"
// @Success 201 {object}
// @Router /v1/auth/refresh [post]
func (h *RefreshHandler) Handle(c *fiber.Ctx) error {
	// TODO: change dto
	dto := &authDtos.RegistrationRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	// TODO: extract user from context
	// TODO: impl refresh token logic here

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       "user_id",
			"access_token":  "access_token",
			"refresh_token": "refresh_token",
		},
	})

	return nil
}
