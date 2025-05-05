//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"layout/internal/biz"
	"layout/internal/conf"
	"layout/internal/data"
	"layout/internal/server"
	"layout/internal/service"
	"layout/pkg/datasource"
	"layout/pkg/monitor"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(context.Context, *conf.Bootstrap, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			datasource.DatasourceProviderSet,
			monitor.MonitorProviderSet,
			server.SrvrProviderSet,
			data.DataProviderSet,
			biz.BizProviderSet,
			service.ServiceProviderSet,
			newApp,
		),
	)
}
