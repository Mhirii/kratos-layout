package data

import (
	"layout/internal/conf"
	"layout/pkg/datasource"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var DataProviderSet = wire.NewSet(NewData, NewUsersRepo, NewProductsRepo)

// dataStruct .
type dataStruct struct {
	gorm        *gorm.DB
	gormCleanup func()

	mongo        *mongo.Database
	mongoCleanup func()

	redis        *redis.Client
	redisCleanup func()

	nats        *nats.Conn
	natsCleanup func()

	js        *jetstream.JetStream
	jsCleanup func()

	log *log.Helper
}

type Data interface {
	GetMongoDB() *mongo.Database
	CleanMongo()

	GetGormDB() *gorm.DB
	CleanGorm()

	GetRedis() *redis.Client
	CleanRedis()

	GetNats() *nats.Conn
	CleanNats()

	GetJetStream() *jetstream.JetStream
	CleanJetStream()
}

func (d *dataStruct) GetMongoDB() *mongo.Database {
	return d.mongo
}

func (d *dataStruct) GetGormDB() *gorm.DB {
	return d.gorm
}

func (d *dataStruct) CleanMongo() {
	d.mongoCleanup()
}

func (d *dataStruct) CleanGorm() {
	d.gormCleanup()
}

func (d *dataStruct) GetRedis() *redis.Client {
	return d.redis
}

func (d *dataStruct) CleanRedis() {
	d.redisCleanup()
}

func (d *dataStruct) GetNats() *nats.Conn {
	return d.nats
}

func (d *dataStruct) CleanNats() {
	d.natsCleanup()
}

func (d *dataStruct) GetJetStream() *jetstream.JetStream {
	return d.js
}

func (d *dataStruct) CleanJetStream() {
	d.jsCleanup()
}

// NewData .
func NewData(
	c *conf.Data,
	g datasource.Gorm,
	m datasource.Mongo,
	n datasource.Nats,
	r datasource.Redis,
	logger log.Logger,
	tp trace.TracerProvider,
) (Data, error) {
	data := &dataStruct{log: log.NewHelper(logger)}
	data.log.Debug("Initializing Data")

	data.log.Debug("Initializing Gorm")
	if g.GetDB() == nil {
		data.log.Warn("No Gorm Database found, skipping gorm initialization")
	} else {
		data.gorm = g.GetDB()
		data.gormCleanup = g.GetCleanup()
		data.log.Debug("Initialized Gorm Successfully! 🎉")
	}

	data.log.Debug("Initializing Mongo")
	if m.GetDB() == nil {
		data.log.Warn("No Mongo Client found, skipping mongo initialization")
	} else {
		data.mongo = m.GetDB()
		data.mongoCleanup = m.GetCleanup()
		data.log.Debug("Initialized Mongo Successfully! 🎉")
	}

	data.log.Debug("Initializing Redis")
	if r.GetClient() == nil {
		data.log.Warn("No Redis client found, skipping redis initialization")
	} else {
		data.redis = r.GetClient()
		data.redisCleanup = r.GetCleanup()
		data.log.Debug("Initialized Redis Successfully! 🎉")
	}

	data.log.Debug("Initializing NATS")
	if n.GetClient() == nil {
		data.log.Warn("No Nats client found, skipping nats")
	} else {
		data.nats = n.GetClient()
		data.natsCleanup = n.GetCleanup()
		data.log.Debug("Initialized NATS Successfully! 🎉")
	}

	data.log.Debug("Initializing JetStream")
	if n.GetJetStream() == nil {
		data.log.Warn("Jetstream is disabled, skipping Jetstream initialization")
	} else {
		data.js = n.GetJetStream()
		data.jsCleanup = n.GetCleanup()
		data.log.Debug("Initialized JetStream Successfully! 🎉")
	}

	return data, nil
}
