package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ServiceProviderSet = wire.NewSet(NewUsersService, NewProductsService)
