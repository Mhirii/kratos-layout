package datasource

import (
	"context"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gtracing "gorm.io/plugin/opentelemetry/tracing"
)

type gormStruct struct {
	db      *gorm.DB
	cleanup func()
	logger  log.Logger
}

type Gorm interface {
	GetDB() *gorm.DB
	GetCleanup() func()
}

func (g *gormStruct) GetDB() *gorm.DB {
	return g.db
}

func (g *gormStruct) GetCleanup() func() {
	return g.cleanup
}

func openDB(ctx context.Context, c *conf.Data, tp trace.TracerProvider) (*gorm.DB, error) {
	driver := c.GetPostgres().GetDriver()
	if driver == "" {
		return nil, errors.InternalServer("no database configured", "no database configured")
	}
	source := c.GetPostgres().GetSource()
	if source == "" {
		return nil, errors.InternalServer("no database configured", "no database configured")
	}
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DriverName:           driver,
				DSN:                  source,
				PreferSimpleProtocol: true,
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}
	err = db.Use(gtracing.NewPlugin(gtracing.WithTracerProvider(tp)))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewGorm(c *conf.Data, logger log.Logger, tp trace.TracerProvider) (Gorm, error) {
	lg := log.NewHelper(logger)

	db, err := openDB(context.Background(), c, tp)
	if err != nil {
		err = errors.InternalServer("Failed to connect to PostgreSQL", err.Error())
		lg.Error(err)
		return nil, err
	}

	cleanup := func() {
	}

	return &gormStruct{
		db:      db,
		logger:  logger,
		cleanup: cleanup,
	}, nil
}

func GormMigrate(ctx context.Context, c *conf.Data, logger log.Logger, models ...interface{}) {
	l := log.NewHelper(logger)
	client, err := openDB(ctx, c, otel.GetTracerProvider())
	if err != nil {
		l.Errorf("failed opening database: %s", err)
	}
	for _, model := range models {
		migrator := client.Migrator()
		err = migrator.AutoMigrate(model)
		if err != nil {
			l.Errorf("failed migrating the schema: %s", err)
			continue
		} else {
			l.Infof("migrated the schema successfully")
		}
		break
	}
	if err != nil {
		l.Errorf("failed migrating the schema: %s", err)
	}
}
