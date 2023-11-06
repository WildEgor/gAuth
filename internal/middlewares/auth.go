package middlewares

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/pkg/errors"
	"strings"

	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

type AuthMiddlewareConfig struct {
	Filter func(c *fiber.Ctx) bool

	UserRepo        *repositories.UserRepository
	KeycloakAdapter *kcAdapter.KeycloakAdapter
	KeycloakConfig  *configs.KeycloakConfig

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

			publicKey, err := parseKeycloakRSAPublicKey(cfg.KeycloakConfig.RSAPublicKey)
			if err != nil {
				return nil, errors.Wrap(err, "parse rsa error")
			}

			parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, errors.New(fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"]))
				}

				return publicKey, nil
			})

			claims, ok := parsedToken.Claims.(jwt.MapClaims)

			res, err := cfg.KeycloakAdapter.UserInfoByToken(context.Background(), token)
			if err != nil {
				return nil, errors.Wrap(err, "invalid token")
			}

			log.Info("[AuthMiddleware] token: %v", token)

			var jwtPayload jwt.MapClaims

			if ok && parsedToken.Valid {

				jwtPayload = jwt.MapClaims{
					"sub":           res.Email,
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

func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
