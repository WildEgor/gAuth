package refresh_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	"github.com/WildEgor/gAuth/internal/configs"
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/gofiber/fiber/v2"
)

type RefreshHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewRefreshHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *RefreshHandler {
	return &RefreshHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
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

	token, err := h.jwt.ParseToken(rt)
	if err != nil {
		resp.SetStatus(c, fiber.StatusBadRequest)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: err.Error(),
		})
		resp.FormResponse()
		return nil
	}

	userId, err := h.tr.GetRT(token.TokenUuid)
	if err != nil {
		resp.SetStatus(c, fiber.StatusBadRequest)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: err.Error(),
		})
		resp.FormResponse()
		return nil
	}

	user, err := h.ur.FindById(userId)
	if err != nil {
		resp.SetStatus(c, fiber.StatusBadRequest)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: err.Error(),
		})
		resp.FormResponse()
		return nil
	}

	nat, atErr := h.jwt.GenerateToken(user.Id.Hex(), h.jwtConfig.ATDuration)
	if atErr != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: atErr.Error(),
		})
		resp.FormResponse()

		return nil
	}

	sErr := h.tr.SetAT(nat)
	if sErr != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: sErr.Error(),
		})
		resp.FormResponse()

		return nil
	}

	//c.Cookie(&fiber.Cookie{
	//	Name:     "access_token",
	//	Value:    *nat.Token,
	//	Path:     "/",
	//	MaxAge:   int(h.jwtConfig.ATDuration.Seconds()),
	//	Secure:   false,
	//	HTTPOnly: true,
	//	Domain:   "localhost",
	//})

	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(fiber.Map{
		"user_id":       user.Id.Hex(),
		"access_token":  nat.Token,
		"refresh_token": rt,
	})
	resp.FormResponse()

	return nil
}
