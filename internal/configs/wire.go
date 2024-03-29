package configs

import "github.com/google/wire"

var ConfigsSet = wire.NewSet(
	NewConfigurator,
	NewAppConfig,
	NewMongoDBConfig,
	NewRedisConfig,
	NewJWTConfig,
	NewOTPConfig,
)
