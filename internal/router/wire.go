package router

import (
	"github.com/WildEgor/gAuth/internal/handlers"
	"github.com/google/wire"
)

var RouterSet = wire.NewSet(
	NewPublicRouter,
	NewSwaggerRouter,
	handlers.HandlersSet,
)
