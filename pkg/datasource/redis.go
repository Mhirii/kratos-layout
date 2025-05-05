package datasource

import (
	"context"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type redisStruct struct {
	log *log.Helper

	client  *redis.Client
	cleanup func()
}

type Redis interface {
	GetClient() *redis.Client
	GetCleanup() func()
}

func (r *redisStruct) GetClient() *redis.Client {
	return r.client
}

func (r *redisStruct) GetCleanup() func() {
	return r.cleanup
}

func NewRedis(c *conf.Data, logger log.Logger, tp trace.TracerProvider) (Redis, error) {
	r := &redisStruct{log: log.NewHelper(logger)}
	if c.GetRedis() == nil {
		err := errors.InternalServer("no redis configured", "no redis configured")
		r.log.Error(err)
		return nil, err
	}
	url := c.GetRedis().GetAddr()
	if url == "" {
		err := errors.InternalServer("no redis configured", "no redis configured")
		r.log.Error(err)
		return nil, err
	}
	network := c.GetRedis().GetNetwork()
	if network == "" {
		r.log.Debug("REDIS: network not found, using tcp")
		network = "tcp"
	}
	options := &redis.Options{
		Addr:    url,
		Network: network,
		DB:      0,
	}

	rc := redis.NewClient(options)
	if rc == nil {
		err := errors.InternalServer("failed to connect to redis", "failed to connect to redis")
		r.log.Error(err)
		return nil, err
	}

	r.log.Debug("REDIS: testing the connection to redis")
	st := rc.Ping(context.Background())
	if st.Err() != nil {
		err := errors.InternalServer("failed to connect to redis", "failed to connect to redis")
		r.log.Error(err)
		return nil, err
	}
	r.log.Debug("REDIS: pinged the connection to redis successfully")

	trOpts := []redisotel.TracingOption{}
	err := redisotel.InstrumentTracing(rc, trOpts...)
	if err != nil {
		err = errors.InternalServer("failed to setup tracing for redis", err.Error())
		r.log.Error(err)
	} else {
		r.log.Debug("REDIS: tracing setup successfully")
	}

	mtOpts := []redisotel.MetricsOption{}
	err = redisotel.InstrumentMetrics(rc, mtOpts...)
	if err != nil {
		err = errors.InternalServer("failed to setup metrics for redis", err.Error())
		r.log.Error(err)
	} else {
		r.log.Debug("REDIS: metrics setup successfully")
	}

	r.cleanup = func() {
		r.log.Warn("REDIS: closing redis connection")
		err := rc.Close()
		if err != nil {
			err = errors.InternalServer("Failed to close redis connection", err.Error())
			r.log.Error(err)
		}
	}
	r.client = rc

	return r, nil
}
