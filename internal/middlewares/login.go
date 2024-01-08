package middlewares

import (
	"errors"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type LoginMiddlewareConfig struct {
	Filter       func(c *fiber.Ctx) bool
	UserRepo     *repositories.UserRepository
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
}

var LoginMiddlewareConfigDefault = LoginMiddlewareConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
}

func configLoginDefault(config ...LoginMiddlewareConfig) LoginMiddlewareConfig {
	if len(config) < 1 {
		return LoginMiddlewareConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = LoginMiddlewareConfigDefault.Filter
	}

	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {
			payload := &authDtos.LoginRequestDto{}
			if err := validators.ParseAndValidate(c, payload); err != nil {
				return nil, errors.New("login/password required")
			}

			// TODO: impl login logic here using login and password

			jwtPayload := jwt.MapClaims{
				"sub":           payload.Login,
				"typ":           "Bearer",
				"exp":           "TODO: impl",
				"access_token":  "TODO: impl",
				"refresh_token": "TODO: impl",
			}

			return &jwtPayload, nil
		}
	}

	// Set default Unauthorized if not passed
	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}

	return cfg
}

// NewLoginMiddleware login via email/phone and password user in Keycloak and extract from DB
func NewLoginMiddleware(config LoginMiddlewareConfig) fiber.Handler {
	// For setting default config
	cfg := configLoginDefault(config)

	return func(c *fiber.Ctx) error {
		payload := &authDtos.LoginRequestDto{}
		if err := validators.ParseAndValidate(c, payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"isOk": false,
				"error": fiber.Map{
					"message": "Bad request",
				},
			})
		}
		// Don't execute middleware if Filter returns true
		if cfg.Filter != nil && cfg.Filter(c) {
			log.Debug("[LoginMiddleware] was skipped")
			return c.Next()
		}
		log.Debug("[LoginMiddleware] was run")

		claims, err := cfg.Decode(c)
		if err == nil {
			c.Locals("jwtClaims", *claims)
			user, err := cfg.UserRepo.FindByLogin(payload.Login, payload.Password)
			if err == nil {
				c.Locals("user", *user)
				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"isOk": false,
			"error": fiber.Map{
				"message": "Unauthorized",
			},
		})
	}
}
