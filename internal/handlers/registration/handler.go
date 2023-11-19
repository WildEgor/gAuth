package registration_handler

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct {
	userRepo  *repositories.UserRepository
	kcAdapter *kcAdapter.KeycloakAdapter
}

func NewRegistrationHandler(
	userRepo *repositories.UserRepository,
	kcAdapter *kcAdapter.KeycloakAdapter,
) *RegistrationHandler {
	return &RegistrationHandler{
		userRepo:  userRepo,
		kcAdapter: kcAdapter,
	}
}

// Handle RegistrationHandler method to create a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 201 {object} authDtos.RegistrationRequestDto
// @Router /v1/user/reg [post]
func (h *RegistrationHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.RegistrationRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	var keycloakUser = gocloak.User{
		Username:      gocloak.StringP(dto.Phone),
		FirstName:     gocloak.StringP(dto.FirstName),
		LastName:      gocloak.StringP(dto.LastName),
		Email:         gocloak.StringP(dto.Email),
		EmailVerified: gocloak.BoolP(true),
		Enabled:       gocloak.BoolP(true),
		Attributes: &map[string][]string{
			"mobile": {dto.Phone},
		},
	}

	_, createErr := h.kcAdapter.CreateUser(context.Background(), keycloakUser, dto.Password, "user")
	if createErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})

		return nil
	}

	userModel := models.UsersModel{
		Email:     dto.Email,
		Phone:     dto.Phone,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}

	newUser, mongoErr := h.userRepo.Create(userModel)
	if mongoErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: err.Error(),
			},
		})

		return nil
	}

	res, loginErr := h.kcAdapter.Login(dto.Email, dto.Password)
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
			"user_id":       newUser.Id.Hex(),
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
		},
	})

	return nil
}
