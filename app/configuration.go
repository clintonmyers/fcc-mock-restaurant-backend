package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Configuration struct {
	DB                  *gorm.DB
	WebApp              *fiber.App
	MaxIdle             int
	MaxOpenConn         int
	LifetimeMinutes     int
	Production          bool
	Port                string
	LocalDB             string
	DatabaseUrl         string
	GenerateTestData    bool
	ApiKey              string
	AutoMigrate         bool
	DeleteLocalDatabase bool
}

const (
	API_KEY_OS      string = "API_KEY"
	API_KEY_FLAG    string = "apiKey"
	API_KEY_DEFAULT string = ""

	DELETE_LOCAL_DB_OS      string = "DELETE_LOCAL_DATABASE"
	DELETE_LOCAL_DB_FLAG    string = "deleteLocalDatabase"
	DELETE_LOCAL_DB_DEFAULT bool   = true

	GENERATE_TEST_DATA_OS      string = "GENERATE_TEST_DATA"
	GENERATE_TEST_DATA_FLAG    string = "generateTestData"
	GENERATE_TEST_DATA_DEFAULT bool   = true

	DATABASE_URL_OS      string = "DATABASE_URL"
	DATABASE_URL_FLAG    string = "databaseUrl"
	DATABASE_URL_DEFAULT string = ""

	LOCAL_DB_OS      string = "LOCAL_DB"
	LOCAL_DB_FLAG    string = "localDB"
	LOCAL_DB_DEFAULT string = "test.db"

	LIFETIME_MINUTES_OS      string = "LIFETIME_MINUTES"
	LIFETIME_MINUTES_FLAG    string = "lifetimeMinutes"
	LIFETIME_MINUTES_DEFAULT int    = 10

	MAX_OPEN_CONN_OS      string = "MAX_OPEN_CONN"
	MAX_OPEN_CONN_FLAG    string = "maxOpenConn"
	MAX_OPEN_CONN_DEFAULT int    = 10

	MAX_IDLE_OS      string = "MAX_IDLE"
	MAX_IDLE_FLAG    string = "maxIdle"
	MAX_IDLE_DEFAULT int    = 5

	PRODUCTION_OS      string = "PRODUCTION"
	PRODUCTION_FLAG    string = "production"
	PRODUCTION_DEFAULT bool   = false

	AUTO_MIGRATE_OS      string = "AUTO_MIGRATE"
	AUTO_MIGRATE_FLAG    string = "autoMigrate"
	AUTO_MIGRATE_DEFAULT bool   = true

	PORT_OS      string = "PORT"
	PORT_FLAG    string = "port"
	PORT_DEFAULT string = ":3030"
)
