package reg_handler

import (
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type RegHandler struct {
	userRepo *repositories.UserRepository
}

func NewRegHandler(
	userRepo *repositories.UserRepository,
) *RegHandler {
	return &RegHandler{
		userRepo: userRepo,
	}
}

// Handle RegHandler method to create a new user
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 201 {object} authDtos.RegistrationRequestDto
// @Router /v1/user/reg [post]
func (h *RegHandler) Handle(c *fiber.Ctx) error {
	dto := &authDtos.RegistrationRequestDto{}
	err := validators.ParseAndValidate(c, dto)
	if err != nil {
		return err
	}

	// TODO

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

	c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       newUser.Id.Hex(),
			"access_token":  "access_token",
			"refresh_token": "refresh_token",
		},
	})

	return nil
}
