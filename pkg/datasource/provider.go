package datasource

import "github.com/google/wire"

var DatasourceProviderSet = wire.NewSet(NewGorm, NewMongo, NewRedis, NewNats)
