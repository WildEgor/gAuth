package reg_handler

import (
	"github.com/WildEgor/gAuth/internal/configs"
	domains "github.com/WildEgor/gAuth/internal/domain"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/mappers"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
)

type RegHandler struct {
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewRegHandler(
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *RegHandler {
	return &RegHandler{
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
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

	// 1. Check if user exists
	existed, _ := h.ur.FindByEmail(dto.Email)
	if existed != nil {
		c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "User already exists",
			},
		})

		return nil
	}

	// 2. Create user of not exists
	userModel := mappers.CreateUser(dto)
	newUser, mongoErr := h.ur.Create(userModel)
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

	// 3. Generate tokens
	at, atErr := h.jwt.GenerateToken(newUser.Id.Hex(), h.jwtConfig.ATDuration)
	rt, rtErr := h.jwt.GenerateToken(newUser.Id.Hex(), h.jwtConfig.ATDuration)
	if atErr != nil || rtErr != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR: tokens", // TODO: make better
			},
		})

		return nil
	}

	errAT := h.tr.SetAT(at)
	errRT := h.tr.SetRT(rt)
	if errAT != nil || errRT != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"isOk": false,
			"data": &domains.ErrorResponseDomain{
				Status:  "fail",
				Message: "ERR", // TODO: make better
			},
		})

		return nil
	}

	// 4. Return tokens
	//c.Cookie(&fiber.Cookie{
	//	Name:     "access_token",
	//	Value:    *at.Token,
	//	Path:     "/",
	//	MaxAge:   int(h.jwtConfig.ATDuration.Seconds()),
	//	Secure:   false,
	//	HTTPOnly: true,
	//	Domain:   "localhost",
	//})
	//
	//c.Cookie(&fiber.Cookie{
	//	Name:     "refresh_token",
	//	Value:    *rt.Token,
	//	Path:     "/",
	//	MaxAge:   int(h.jwtConfig.RTDuration.Seconds()),
	//	Secure:   false,
	//	HTTPOnly: true,
	//	Domain:   "localhost",
	//})

	// 5. Or response with tokens
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"isOk": true,
		"data": fiber.Map{
			"user_id":       newUser.Id.Hex(),
			"access_token":  at.Token,
			"refresh_token": rt.Token,
		},
	})
}
