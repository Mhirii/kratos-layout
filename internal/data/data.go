package data

import (
	"github.com/go-kratos/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData, NewUsersRepo)

type Data struct {
	// TODO wrapped database client
	// client *gorm.DB
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

/*

var driver = os.Getenv("DRIVER")
var host = os.Getenv("DB_HOST")
var name = os.Getenv("DB_NAME")
var port = os.Getenv("DB_PORT")
var user = os.Getenv("DB_USER")
var pass = os.Getenv("DB_PASS")
var sslmode = os.Getenv("DB_SSLMODE")


func verifyEnv(logger log.Logger) bool {
	for _, v := range []string{driver, host, name, port, user, pass} {
		if v == "" {
			log.NewHelper(logger).Error("missing environment variable %s", v)
			return false
		}
	}
	if sslmode == "" {
		log.NewHelper(logger).Warn("missing environment variable DB_SSLMODE, using default value 'disable'")
		sslmode = "disable"
	}
	return true
}

func getDSN() string {
	return "host=" + host + " port=" + port + " user=" + user + " dbname=" + name + " password=" + pass + " sslmode=" + sslmode
}

func openDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  getDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	if !verifyEnv(logger) {
		return nil, nil, errors.New("missing environment variable")
	}
	client, err := openDB()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		//err := client.Close()
		if err != nil {
			log.NewHelper(logger).Error("failed closing the data resources: %v", err)
		}
	}
	return &Data{client}, cleanup, nil
}

func Migrate() {
	log.NewHelper(log.DefaultLogger).Info("migrating the schema")
	client, err := openDB()
	if err != nil {
		log.NewHelper(log.DefaultLogger).Error("failed opening database: %v", err)
	}
	err = client.AutoMigrate(&Users{})
	if err != nil {
		log.NewHelper(log.DefaultLogger).Error("failed migrating the schema: %v", err)
	}
}

*/
