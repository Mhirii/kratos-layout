package datasource

import (
	"context"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.opentelemetry.io/otel/trace"
)

type mongoStruct struct {
	db      *mongo.Database
	log     log.Logger
	cleanup func()
}

type Mongo interface {
	GetDB() *mongo.Database
	GetCleanup() func()
}

func (m *mongoStruct) GetDB() *mongo.Database {
	return m.db
}

func (m *mongoStruct) GetCleanup() func() {
	return m.cleanup
}

func connMongo(ctx context.Context, c *conf.Data, logger log.Logger, database string, tp trace.TracerProvider) (*mongo.Database, *mongo.Client, error) {
	lg := log.NewHelper(logger)
	uri := c.GetMongo().GetUri()
	if uri == "" {
		lg.Warn("MONGO: mongodb uri is empty, using default value mongodb://localhost:27017")
		uri = "mongodb://localhost:27017"
	}
	opts := options.Client()
	opts.ApplyURI(uri)

	cm := otelmongo.NewMonitor(otelmongo.WithTracerProvider(tp), otelmongo.WithCommandAttributeDisabled(true))
	opts.SetMonitor(cm)

	username := c.GetMongo().GetUsername()
	password := c.GetMongo().GetPassword()

	if username != "" && password != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		lg.Error("failed to connect to mongodb", err)
		return nil, nil, err
	}

	lg.Info("MONGO: pinging mongodb")
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		lg.Error("failed to ping mongodb", err)
		panic(err)
	}

	db := client.Database(database)
	return db, client, nil
}

func NewMongo(ctx context.Context, c *conf.Data, logger log.Logger, tp trace.TracerProvider) (Mongo, error) {
	lg := log.NewHelper(logger)
	lg.Info("MONGO: Initiating NewData")

	database := c.GetMongo().GetDatabase()
	db, client, err := connMongo(ctx, c, logger, database, tp)
	if err != nil {
		return nil, err
	}

	cleanup := func() {
		lg.Info("MONGO: closing the data resources")
		err := client.Disconnect(context.Background())
		if err != nil {
			lg.Error("Error disconnecting mongodb: ", err)
		}
		lg.Info("MONGO: disconnected from mongodb successfully")
	}

	log.NewHelper(logger).Info("connected to database successfully")
	return &mongoStruct{
		db:      db,
		log:     logger,
		cleanup: cleanup,
	}, nil
}
