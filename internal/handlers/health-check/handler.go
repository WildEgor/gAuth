package health_check_handler

import (
	core_dtos "github.com/WildEgor/g-core/pkg/core/dtos"
	"github.com/WildEgor/gAuth/internal/configs"
	domains "github.com/WildEgor/gAuth/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type HealthCheckHandler struct {
	appConfig *configs.AppConfig
}

func NewHealthCheckHandler(
	appConfig *configs.AppConfig,
) *HealthCheckHandler {
	return &HealthCheckHandler{
		appConfig,
	}
}

func (hch *HealthCheckHandler) Handle(c *fiber.Ctx) error {
	resp := core_dtos.InitResponse()
	resp.SetStatus(c, fiber.StatusOK)
	resp.SetData(&domains.StatusDomain{
		Status:      "ok",
		Version:     hch.appConfig.Version,
		Environment: hch.appConfig.GoEnv,
	})
	resp.FormResponse()
	return nil
}
