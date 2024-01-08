package middlewares

import (
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type AuthMiddlewareConfig struct {
	Filter       func(c *fiber.Ctx) bool
	UR           *repositories.UserRepository
	JWT          *services.JWTAuthenticator
	Unauthorized fiber.Handler
	Decode       func(c *fiber.Ctx) (*jwt.MapClaims, error)
}

var AuthMiddlewareConfigDefault = AuthMiddlewareConfig{
	Filter:       nil,
	Decode:       nil,
	Unauthorized: nil,
}

func configAuthDefault(config ...AuthMiddlewareConfig) AuthMiddlewareConfig {
	if len(config) < 1 {
		return AuthMiddlewareConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = LoginMiddlewareConfigDefault.Filter
	}

	if cfg.Decode == nil {
		// Set default Decode function if not passed
		cfg.Decode = func(c *fiber.Ctx) (*jwt.MapClaims, error) {
			var token string
			authHeader := c.Get("Authorization")

			if len(authHeader) > 0 {
				components := strings.SplitN(authHeader, " ", 2)

				if len(components) == 2 && components[0] == "Bearer" {
					token = components[1]
				}
			}

			if len(token) == 0 {
				return nil, errors.New("empty token")
			}

			// TODO: need check token in Redis too
			claims, err := cfg.JWT.ParseToken(token)
			jwtPayload := jwt.MapClaims{}

			if err == nil && claims != nil && claims.IsValid {
				jwtPayload = jwt.MapClaims{
					"sub":           claims.UserID,
					"typ":           "Bearer",
					"exp":           claims.ExpiresIn,
					"access_token":  authHeader,
					"refresh_token": "",
				}

				c.Locals("access_token_uuid", claims.TokenUuid)

			} else {
				return nil, errors.Wrap(err, "token validation")
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

// NewAuthMiddleware validate accessToken in Keycloak and parse it, extract user from DB
func NewAuthMiddleware(config AuthMiddlewareConfig) fiber.Handler {
	// For setting default config
	cfg := configAuthDefault(config)

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Filter returns true
		if cfg.Filter != nil && cfg.Filter(c) {
			log.Debug("[AuthMiddleware] was skipped")
			return c.Next()
		}
		log.Debug("[AuthMiddleware] was run")

		claims, err := cfg.Decode(c)
		if err == nil {
			c.Locals("jwtClaims", *claims)

			id := (*claims)["sub"].(string)
			user, err := cfg.UR.FindById(id)
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
