package proto

import (
	"github.com/google/wire"
)

var GRPCSet = wire.NewSet(
	NewGRPCServer,
)
