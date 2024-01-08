package login_handler

import (
	"fmt"
	"github.com/WildEgor/gAuth/internal/configs"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type LoginHandler struct {
	ur        *repositories.UserRepository
	jwtConfig *configs.JWTConfig
}

func NewLoginHandler(
	ur *repositories.UserRepository,
	jwtConfig *configs.JWTConfig,
) *LoginHandler {
	return &LoginHandler{
		ur:        ur,
		jwtConfig: jwtConfig,
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

	// 4. Return tokens
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    fmt.Sprintf("%s", jwtClaims["access_token"]),
		Path:     "/",
		MaxAge:   int(h.jwtConfig.ATDuration.Seconds()),
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    fmt.Sprintf("%s", jwtClaims["refresh_token"]),
		Path:     "/",
		MaxAge:   int(h.jwtConfig.RTDuration.Seconds()),
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       authUser.Id.Hex(),
			"access_token":  jwtClaims["access_token"],
			"refresh_token": jwtClaims["refresh_token"],
		},
	})

	return nil
}
