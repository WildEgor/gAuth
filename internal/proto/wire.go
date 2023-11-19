package proto

import (
	"github.com/google/wire"
)

var RPCSet = wire.NewSet(
	NewAuthService,
	NewProxyService,
	NewGRPCServer,
)
