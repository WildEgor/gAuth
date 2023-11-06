package login_handler

import (
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type LoginHandler struct {
	userRepo *repositories.UserRepository
}

func NewLoginHandler(
	userRepo *repositories.UserRepository,
) *LoginHandler {
	return &LoginHandler{
		userRepo: userRepo,
	}
}

// Handle LoginHandler method login via email/password or phone/password
// @Description Login user
// @Summary login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body string true "Login"
// @Param password body string true "Password"
// @Success 200 {object} authDtos.LoginRequestDto
// @Router /v1/auth/login [post]
func (h *LoginHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.LoginRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	jwtClaims := c.Locals("jwtClaims").(jwt.MapClaims)
	if jwtClaims == nil {
		c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"isOk": true,
			"data": fiber.Map{
				"message": "Not authorized",
			},
		})

		return nil
	}

	authUser := middlewares.ExtractUser(c)

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       authUser.Id.Hex(),
			"access_token":  jwtClaims["access_token"],
			"refresh_token": jwtClaims["refresh_token"],
		},
	})

	return nil
}
