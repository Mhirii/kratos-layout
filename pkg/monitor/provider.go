package monitor

import "github.com/google/wire"

var MonitorProviderSet = wire.NewSet(
	NewMeter,
	NewMeterProvider,
	NewTracer,
	NewTracerProvider,
	NewTextMapPropagator,
)
