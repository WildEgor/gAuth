package change_password_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	"github.com/WildEgor/gAuth/internal/configs"
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type ChangePasswordHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewChangePasswordHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

// Handle ChangePasswordHandler method allows change password
// @Description Allow change authenticated user own password
// @Summary change password
// @Tags Auth
// @Accept json
// @Produce json
// @Param old_password body string true "OldPassword"
// @Param new_password body string true "NewPassword"
// @Success 200 {object} authDtos.ChangePasswordRequestDto
// @Router /v1/auth/change-password [post]
func (h *ChangePasswordHandler) Handle(c *fiber.Ctx) error {
	dto := &user.ChangePasswordRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	resp := core_dtos.InitResponse()

	authUser := middlewares.ExtractUser(c)

	_, cmpErr := authUser.ComparePassword(dto.OldPassword)
	if cmpErr != nil {
		resp.SetStatus(c, fiber.StatusOK)
		resp.SetData(fiber.Map{
			"message": "Invalid password",
		})

		resp.FormResponse()

		return cmpErr
	}

	setPassErr := authUser.SetPassword(dto.NewPassword)
	if setPassErr != nil {
		return setPassErr
	}

	upErr := h.ur.UpdatePassword(*authUser)
	if upErr != nil {
		return upErr
	}

	// 3. Generate tokens
	at, atErr := h.jwt.GenerateToken(authUser.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(authUser.Id.Hex(), h.jwtConfig.ATDuration)
	if atErr != nil || rtErr != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR: tokens", // TODO: make better
		})

		resp.FormResponse()

		return nil
	}

	errAT := h.tr.SetAT(at)
	errRT := h.tr.SetRT(rt)
	if errAT != nil || errRT != nil {
		resp.SetStatus(c, fiber.StatusInternalServerError)
		resp.SetData(&domains.ErrorResponseDomain{
			Status:  "fail",
			Message: "ERR", // TODO: make better
		})

		resp.FormResponse()

		return nil
	}

	resp.SetStatus(c, fiber.StatusCreated)
	resp.SetData(fiber.Map{
		"user_id":       authUser.Id.Hex(),
		"access_token":  at.Token,
		"refresh_token": rt.Token,
	})

	resp.FormResponse()

	return nil
}
