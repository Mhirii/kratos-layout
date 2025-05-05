package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
var SrvrProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
