package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"gorm.io/driver/postgres"
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

func setDatabaseParameters(config *app.Configuration) {
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
	// AutoMigrate
	if config.AutoMigrate = strings.ToLower(os.Getenv("AUTO_MIGRATE")) == "true"; config.AutoMigrate == false {
		flag.BoolVar(&config.AutoMigrate, "autoMigrate", true, "Should we auto migrate the database?")
	}

	// Production
	if config.Production = strings.ToLower(os.Getenv("production")) == "true"; config.Production == false {
		flag.BoolVar(&config.Production, "production", false, "Is this a production run?")
	}
	// Production
	if config.GenerateLocalData = strings.ToLower(os.Getenv("GENERATE_LOCAL_DATA")) == "true"; config.GenerateLocalData == false {
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

	if config.DeleteLocalDatabase = strings.ToLower(os.Getenv("DELETE_LOCAL_DATABASE")) == "true"; config.DeleteLocalDatabase == false {
		flag.BoolVar(&config.DeleteLocalDatabase, "deleteLocalDatabase", true, "Delete Local Database upon start")
	}

	// TEMP CONFIGS
	// These should be removed when we're actually doing real work

	if config.API_KEY = os.Getenv("API_KEY"); config.Port == "" {
		flag.StringVar(&config.API_KEY, "API_KEY", "", "API key to use, will default to an empty string which will always deny")
	}
	//
	flag.Parse()

	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	if config.Production == false && config.DeleteLocalDatabase {
		localDatabaseName := "./" + config.LocalDB
		if f, err := os.Stat(localDatabaseName); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			fmt.Println("File does not exist")
		} else {
			fmt.Println("File exists")
			err := os.Remove(localDatabaseName)
			if err != nil {
				//log.Fatal(err)
				fmt.Println(err)
			}
			os.Create(localDatabaseName)
			fmt.Printf("%s -> %d", f.Name(), f.Size())

		}
	}

	// Setup DB connection
	if config.Production {

		//fmt.Println("Connecting to production database")
		//
		//db, err := gorm.
		//	Open(mysql.Open("fcc:Password123@tcp(localhost:3306)/testing?charset=utf8mb4&parseTime=True&loc=Local"),
		//		&gorm.Config{})
		//if err != nil {
		//	panic("failed to connect database")
		//}
		//config.DB = db

		fmt.Println("Connecting to production database")
		if db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{}); err != nil {
			panic("failed to connect database")
		} else {
			config.DB = db
		}

	} else {
		fmt.Println("Connecting to non-production database")
		//if db, err := gorm.Open(sqlite.Open(config.LocalDB), &gorm.Config{
		//	PrepareStmt: true,
		//}); err != nil {
		//	panic("failed to connect database")
		//
		//} else {
		//	config.DB = db
		//}
	}

	// Setup the web app
	config.WebApp = fiber.New()

	fmt.Printf("ENV production: %s\n", os.Getenv("production"))

	if !config.Production {
		fmt.Printf("ENV full configuration: %#v\n", config)
	}

}
