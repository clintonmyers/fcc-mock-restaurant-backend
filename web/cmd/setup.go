package main

import (
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func configureMiddleware(app *fiber.App) {
	// https://github.com/gofiber/fiber/tree/master/middleware/pprof
	// pprof logging is underneath /debug/pprof
	app.Use(pprof.New())

	// https://docs.gofiber.io/api/middleware/cors
	// The default already has '*' as the allowed origins
	app.Use(cors.New())

}

func configureDatabase(config *app.Configuration) {
	// Gets the underlying DB connection

	sqlDb, err := config.DB.DB()

	if err != nil {
		log.Fatal(err)
	}

	// configure the connection limits
	sqlDb.SetMaxIdleConns(config.MaxIdle)
	sqlDb.SetMaxOpenConns(config.MaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(config.LifetimeMinutes))
}

func loadConfiguration(config *app.Configuration) {
	// Port
	if config.Port = os.Getenv("PORT"); config.Port == "" {
		flag.StringVar(&config.Port, "port", ":3000", "Port to use")
	}

	// Production
	if config.Production = strings.ToLower(os.Getenv("production")) == "true"; config.Production == false {
		flag.BoolVar(&config.Production, "production", false, "Is this a production run?")
	}

	// MaxIdle
	if s, err := strconv.Atoi(os.Getenv("maxIdle")); err == nil {
		config.MaxIdle = s
	} else {
		flag.IntVar(&config.MaxIdle, "maxIdle", 5, "Set max idle connections for database")
	}

	// MaxOpenConn
	if s, err := strconv.Atoi(os.Getenv("maxOpenConn")); err == nil {
		config.MaxOpenConn = s
	} else {
		flag.IntVar(&config.MaxOpenConn, "maxOpenConn", 10, "Set max open connections for database")
	}

	// LifetimeMinutes
	if s, err := strconv.Atoi(os.Getenv("lifetimeMinutes")); err == nil {
		config.LifetimeMinutes = s
	} else {
		flag.IntVar(&config.LifetimeMinutes, "lifetimeMinutes", 10, "Set max open connections for database")
	}

	// DB
	if config.LocalDB = os.Getenv("localDB"); config.LocalDB == "" {
		flag.StringVar(&config.LocalDB, "localDB", "test.db", "Local database file used only during non-production")
	}
	// DatabaseUrl
	if config.DatabaseUrl = os.Getenv("DATABASE_URL"); config.DatabaseUrl == "" {
		fmt.Println("DB URL: ", os.Getenv("DATABASE_URL"))
		flag.StringVar(&config.DatabaseUrl, "databaseUrl", "DEFAULT", "Database URL")
	}

	if config.GenerateTestData = strings.ToLower(os.Getenv("GENERATE_TEST_DATA")) == "true"; config.GenerateTestData == false {
		flag.BoolVar(&config.GenerateTestData, "generateTestData", true, "Generate test data")
	}

	flag.Parse()

	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	// Setup DB connection
	if config.Production {
		fmt.Println("Connecting to production database")
		db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	} else {
		fmt.Println("Connecting to non-production database")
		db, err := gorm.Open(sqlite.Open(config.LocalDB), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	}

	// Setup the web app
	config.WebApp = fiber.New()

	fmt.Printf("ENV production: %s\n", os.Getenv("production"))

	if !config.Production {
		fmt.Printf("ENV full configuration: %#v\n", config)
	}

}
