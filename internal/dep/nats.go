package dep

import (
	"context"
	"time"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

type Nats struct {
	UsesJS  bool
	Conn    *nats.Conn
	JS      jetstream.JetStream
	Logger  *log.Helper
	tp      trace.Tracer
	Cleanup func()
}

func NewNats(c *conf.Bootstrap, logger log.Logger, tp trace.TracerProvider) (*Nats, error) {
	t := tp.Tracer("nats")
	log := log.NewHelper(logger)
	var conn *nats.Conn

	conn, clean, err := connect(c, log, t, context.Background())
	if err != nil {
		return nil, err
	}

	nats := Nats{
		UsesJS:  false,
		Logger:  log,
		tp:      t,
		Conn:    conn,
		Cleanup: clean,
	}

	usesJS := c.GetData().GetNats().GetJetstream()

	if usesJS {
		js, clean, err := nats.connectJS(c)
		if err != nil {
			return nil, err
		}
		return &Nats{
			UsesJS:  usesJS,
			Logger:  log,
			tp:      t,
			Conn:    conn,
			JS:      js,
			Cleanup: clean,
		}, nil
	}

	return &nats, nil
}

func connect(c *conf.Bootstrap, log *log.Helper, t trace.Tracer, ctx context.Context) (*nats.Conn, func(), error) {
	ctx, span := t.Start(ctx, "dep.Nats.conn")
	defer span.End()

	addr := nats.DefaultURL

	if c.GetData().GetNats().GetAddr() == "" {
		log.Warn("No data.nats.addr was set in the config, using default")
	} else {
		addr = c.GetData().GetNats().GetAddr()
	}

	name := "unnamed"
	if c.GetMetadata().GetName() != "" {
		name = c.GetMetadata().GetName()
	}
	opts := []nats.Option{
		nats.Name(name),
		nats.MaxReconnects(-1), // Unlimited reconn
		nats.ReconnectWait(time.Second * 5),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Errorf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("NATS reconnected to %s", nc.ConnectedUrl())
		}),
	}

	cleanup := func() {}

	nc, err := nats.Connect(addr, opts...)
	if err != nil {
		log.Errorf("failed to connect to NATS at %s: %v", addr, err)
		return nil, cleanup, errors.InternalServer("failed to connect to NATS", err.Error())
	}

	cleanup = func() {
		if nc != nil {
			nc.Close()
		}
	}

	return nc, cleanup, nil
}

func (n *Nats) connectJS(c *conf.Bootstrap) (jetstream.JetStream, func(), error) {
	ctx, span := n.tp.Start(context.Background(), "dep.Nats.connJS")
	defer span.End()

	cleanup := func() {}
	var conn *nats.Conn = n.Conn
	var err error

	if conn == nil {
		conn, _, err = connect(c, n.Logger, n.tp, ctx)
		if err != nil {
			return nil, cleanup, err
		}
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, cleanup, err
	}

	cleanup = func() {
		if js.Conn() != nil {
			js.Conn().Close()
		}
	}

	return js, cleanup, nil
}
