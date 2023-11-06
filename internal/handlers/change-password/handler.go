package change_password_handler

import (
	"context"
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/WildEgor/gAuth/internal/dtos/user"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type ChangePasswordHandler struct {
	ka *kcAdapter.KeycloakAdapter
	ur *repositories.UserRepository
}

func NewChangePasswordHandler(
	ka *kcAdapter.KeycloakAdapter,
	ur *repositories.UserRepository,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		ka,
		ur,
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

	authUser := middlewares.ExtractUser(c)

	_, cmpErr := authUser.ComparePassword(dto.OldPassword)
	if cmpErr != nil {
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"isOk": false,
			"data": fiber.Map{
				"message": "Invalid password",
			},
		})

		return cmpErr
	}

	setPassErr := authUser.SetPassword(dto.NewPassword)
	if setPassErr != nil {
		return setPassErr
	}

	setErr := h.ka.UpdatePassword(context.Background(), authUser.Email, dto.NewPassword)
	if setErr != nil {
		return setErr
	}

	upErr := h.ur.Update(*authUser)
	if upErr != nil {
		return upErr
	}

	res, loginErr := h.ka.Login(authUser.Email, dto.NewPassword)
	if loginErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})

		return nil
	}

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       authUser.Id.Hex(),
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
		},
	})

	return nil
}
