package middlewares

import (
	"fmt"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type AuthMiddlewareConfig struct {
	Filter       func(c *fiber.Ctx) bool
	UserRepo     *repositories.UserRepository
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
			headerValue := c.Get("authorization")

			if len(headerValue) > 0 {
				components := strings.SplitN(headerValue, " ", 2)

				if len(components) == 2 && components[0] == "Bearer" {
					token = components[1]
				}
			}

			if len(token) == 0 {
				return nil, errors.New("empty token")
			}

			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
				}

				return "TODO PUB KEY", nil
			})

			claims, ok := parsedToken.Claims.(jwt.MapClaims)

			// TODO: get user from DB

			log.Info("[AuthMiddleware] token: %v", token)

			var jwtPayload jwt.MapClaims

			if ok && parsedToken.Valid {

				jwtPayload = jwt.MapClaims{
					"sub":           "TODO USER ID",
					"typ":           claims["typ"],
					"exp":           claims["exp"],
					"access_token":  token,
					"refresh_token": "",
				}

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

			email := (*claims)["sub"].(string)
			user, err := cfg.UserRepo.FindByEmail(email)
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
