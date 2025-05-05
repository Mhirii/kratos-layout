package datasource

import (
	"time"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
)

type natsStruct struct {
	log     *log.Helper
	cleanup func()
	client  *nats.Conn
	js      *jetstream.JetStream
}

func (n *natsStruct) GetClient() *nats.Conn {
	return n.client
}

func (n *natsStruct) GetCleanup() func() {
	return n.cleanup
}

func (n *natsStruct) GetJetStream() *jetstream.JetStream {
	return n.js
}

type Nats interface {
	GetClient() *nats.Conn
	GetCleanup() func()
	GetJetStream() *jetstream.JetStream
}

func NewNats(c *conf.Data, logger log.Logger, tp trace.TracerProvider) (Nats, error) {
	n := &natsStruct{log: log.NewHelper(logger)}
	if c.GetNats() == nil {
		err := errors.InternalServer("no nats configured", "no nats configured")
		n.log.Error(err)
		return nil, err
	}
	url := c.GetNats().GetAddr()
	if url == "" {
		err := errors.InternalServer("no nats configured", "no nats configured")
		n.log.Error(err)
		return nil, err
	}

	name := c.GetNats().GetName()
	if name == "" {
		name = "kratos"
	}
	options := []nats.Option{
		nats.Name(name),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
	}

	user := c.GetNats().GetUsername()
	pw := c.GetNats().GetPassword()
	if user != "" && pw != "" {
		n.log.Debug("NATS: nats user and password found, adding to options")
		options = append(options, nats.UserInfo(user, pw))
	}

	nc, err := nats.Connect(url, options...)
	if err != nil {
		err = errors.InternalServer("Failed to connect to NATS", err.Error())
		n.log.Error(err)
		return nil, err
	}
	n.client = nc

	if c.GetNats().GetJetstream() {
		n.log.Debug("NATS: connecting to nats jetstream")
		js, err := jetstream.New(nc)
		if err != nil {
			err = errors.InternalServer("Failed to connect to NATS JetStream", err.Error())
			n.log.Error(err)
		}
		if js != nil {
			n.log.Debug("NATS: connected to nats jetstream")
		}
		n.js = &js
	}

	n.cleanup = func() {
		n.log.Warn("NATS: closing nats connection")
		err = nc.Drain()
		if err != nil {
			err = errors.InternalServer("Failed to close nats connection", err.Error())
			n.log.Error(err)
		}
	}

	return n, nil
}
