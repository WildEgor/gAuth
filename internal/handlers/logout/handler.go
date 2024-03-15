package logout_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
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
	resp := core_dtos.InitResponse()

	rt := c.Cookies("refresh_token")
	if rt == "" {
		resp.SetStatus(c, fiber.StatusBadRequest)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "Refresh token is required",
		})
		resp.FormResponse()

		return nil
	}

	token, err := h.jwt.ParseToken(string(rt[:]))
	if err != nil {
		resp.SetStatus(c, fiber.StatusBadRequest)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: err.Error(),
		})
		resp.FormResponse()

		return nil
	}

	atUid := c.Locals("access_token_uuid").(string)
	dErr := h.tr.DeleteTokens(token.TokenUuid, atUid)
	if dErr != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: dErr.Error(),
		})
		resp.FormResponse()

		return nil
	}

	h.resetCookies(c)

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(domains.NewVoidResponseDomain())
	resp.FormResponse()

	return nil
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
