package middlewares

import (
	"errors"

	keycloak_adapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type AuthMiddlewareConfig struct {
	Filter func(c *fiber.Ctx) bool

	KeycloakAdapter *keycloak_adapter.KeycloakAdapter

	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
}

var AuthMiddlewareConfigDefault = AuthMiddlewareConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
}

func configDefault(config ...AuthMiddlewareConfig) AuthMiddlewareConfig {
	if len(config) < 1 {
		return AuthMiddlewareConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = AuthMiddlewareConfigDefault.Filter
	}

	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {

			payload := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{}

			if err := c.BodyParser(&payload); err != nil {
				return nil, errors.New("Email/Password required!")
			}

			res, err := cfg.KeycloakAdapter.Login(payload.Email, payload.Password)
			if err != nil {
				return nil, errors.New("Invalid token")
			}

			log.Info("[AuthMiddleware] token: %v", res.AccessToken)

			// TODO: check validation

			var jwtPayload jwt.MapClaims

			jwtPayload["sub"] = res.IDToken
			jwtPayload["typ"] = res.TokenType
			jwtPayload["token"] = res.AccessToken

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

func NewAuthMiddleware(config AuthMiddlewareConfig) fiber.Handler {
	// For setting default config
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Filter returns true
		if cfg.Filter != nil && cfg.Filter(c) {
			log.Info("[AuthMiddleware] was skipped")
			return c.Next()
		}
		log.Info("[AuthMiddleware] was run")

		claims, err := cfg.Decode(c)

		if err == nil {
			c.Locals("jwtClaims", *claims)
			return c.Next()
		}

		return cfg.Unauthorized(c)
	}
}
