package logout_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/gofiber/fiber/v2"
	"time"
)

type LogoutHandler struct {
	tr  *repositories.TokensRepository
	jwt *services.JWTAuthenticator
}

func NewLogoutHandler(
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
) *LogoutHandler {
	return &LogoutHandler{
		tr:  tr,
		jwt: jwt,
	}
}

// Handle LogoutHandler method logout user
// @Description Logout user
// @Summary Logout user
// @Tags Auth
// @Accept json
// @Produce json
// @Router /v1/auth/logout [post]
func (h *LogoutHandler) Handle(c *fiber.Ctx) error {
	rt := c.Cookies("refresh_token")
	if rt == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "Refresh token is required",
			},
		})

		return nil
	}

	token, err := h.jwt.ParseToken(string(rt[:]))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})

		return nil
	}

	atUid := c.Locals("access_token_uuid").(string)
	dErr := h.tr.DeleteTokens(token.TokenUuid, atUid)
	if dErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: dErr.Error(),
			},
		})

		return nil
	}

	h.resetCookies(c)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"isOk": true,
		"data": domains.NewVoidResponseDomain(),
	})
}

func (h *LogoutHandler) resetCookies(c *fiber.Ctx) {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
	c.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})
}
