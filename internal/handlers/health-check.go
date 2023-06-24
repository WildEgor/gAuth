package handlers

import (
	"github.com/WildEgor/gAuth/internal/config"
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckHandler struct {
	cfg *config.AppConfig
}

func NewHealthCheckHandler(
	cfg *config.AppConfig,
) *HealthCheckHandler {
	return &HealthCheckHandler{
		cfg,
	}
}

func (hch *HealthCheckHandler) Handle(c *fiber.Ctx) error {
	c.JSON(fiber.Map{
		"isOk": true,
		"data": &domains.StatusDomain{
			Status:      "ok",
			Version:     hch.cfg.Version,
			Environment: hch.cfg.GoEnv,
		},
	})
	return nil
}
