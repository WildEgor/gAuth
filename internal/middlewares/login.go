package middlewares

import (
	"errors"
	"github.com/WildEgor/gAuth/internal/configs"
	authDtos "github.com/WildEgor/gAuth/internal/dtos/auth"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/WildEgor/gAuth/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type LoginMiddlewareConfig struct {
	Filter       func(c *fiber.Ctx) bool
	UR           *repositories.UserRepository
	TR           *repositories.TokensRepository
	JWT          *services.JWTAuthenticator
	JWTConfig    *configs.JWTConfig
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

			ur, err := cfg.UR.FindByLogin(payload.Login, payload.Password)
			if err != nil {
				return nil, errors.New("cannot find user")
			}

			at, atErr := cfg.JWT.GenerateToken(ur.Id.Hex(), cfg.JWTConfig.ATDuration)
			if atErr != nil {
				return nil, errors.New("cannot generate token")
			}

			rt, rtErr := cfg.JWT.GenerateToken(ur.Id.Hex(), cfg.JWTConfig.ATDuration)
			if rtErr != nil {
				return nil, errors.New("cannot generate token")
			}

			errAT := cfg.TR.SetAT(at)
			if errAT != nil {
				return nil, errors.New("cannot get token")
			}
			errRT := cfg.TR.SetRT(rt)
			if errRT != nil {
				return nil, errors.New("cannot get token")
			}

			jwtPayload := jwt.MapClaims{
				"sub":           ur.Id.Hex(),
				"typ":           "Bearer",
				"exp":           at.ExpiresIn,
				"access_token":  at.Token,
				"refresh_token": rt.Token,
			}

			c.Locals("access_token_uuid", at.TokenUuid)

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
			user, err := cfg.UR.FindByLogin(payload.Login, payload.Password)
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
