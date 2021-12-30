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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func loadConfiguration(config *app.Configuration) {
	updatePort(config)
	updateAutoMigrate(config)
	updateProductionSetting(config)
	updateGenerateLocalData(config)
	updateMaxIdle(config)
	updateMaxOpenConnections(config)
	updateLifetimeMinutes(config)
	updateLocalDBSetting(config)
	updateDatabaseURL(config)
	updateGenerateTestData(config)
	updateDeleteLocalDatabase(config)

	// TEMP CONFIGS: these should be removed when we're actually doing real authentication
	updateAPIKey(config)

	flag.Parse()

	// the port requires a colon in front of it
	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	if err := localDBFileMaintenance(config); err != nil {
		log.Fatal(err)
	}

	if err := connectToDatabase(config); err != nil {
		log.Fatal(err)
	}

	config.WebApp = fiber.New()

	if !config.Production {
		fmt.Printf("ENV full configuration: %#v\n", config)

	} else {
		fmt.Printf("ENV production: %s\n", os.Getenv("production"))
	}

}

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
	if sqlDb, err := config.DB.DB(); err == nil {
		// configure the connection limits
		sqlDb.SetMaxIdleConns(config.MaxIdle)
		sqlDb.SetMaxOpenConns(config.MaxOpenConn)
		sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(config.LifetimeMinutes))
	} else {
		log.Fatal(err)
	}
}

func updatePort(config *app.Configuration) {
	if config.Port = os.Getenv("PORT"); config.Port == "" {
		flag.StringVar(&config.Port, "port", ":3030", "Port to use")
	}
}
func updateAutoMigrate(config *app.Configuration) {
	if config.AutoMigrate = strings.ToLower(os.Getenv("AUTO_MIGRATE")) == "true"; config.AutoMigrate == false {
		flag.BoolVar(&config.AutoMigrate, "autoMigrate", true, "Should we auto migrate the database?")
	}
}

func updateProductionSetting(config *app.Configuration) {
	if config.Production = strings.ToLower(os.Getenv("production")) == "true"; config.Production == false {
		flag.BoolVar(&config.Production, "production", false, "Is this a production run?")
	}
}
func updateGenerateLocalData(config *app.Configuration) {
	if config.GenerateLocalData = strings.ToLower(os.Getenv("GENERATE_LOCAL_DATA")) == "true"; config.GenerateLocalData == false {
		flag.BoolVar(&config.GenerateLocalData, "generateData", false, "Do we generate local data?")
	}
}

func updateMaxIdle(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv("maxIdle")); err == nil {
		config.MaxIdle = s
	} else {
		flag.IntVar(&config.MaxIdle, "maxIdle", 5, "Set max idle connections for database")
	}
}
func updateMaxOpenConnections(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv("maxOpenConn")); err == nil {
		config.MaxOpenConn = s
	} else {
		flag.IntVar(&config.MaxOpenConn, "maxOpenConn", 10, "Set max open connections for database")
	}
}
func updateLifetimeMinutes(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv("lifetimeMinutes")); err == nil {
		config.LifetimeMinutes = s
	} else {
		flag.IntVar(&config.LifetimeMinutes, "lifetimeMinutes", 10, "Set max open connections for database")
	}
}
func updateLocalDBSetting(config *app.Configuration) {
	if config.LocalDB = os.Getenv("localDB"); config.LocalDB == "" {
		flag.StringVar(&config.LocalDB, "localDB", "test.db", "Local database file used only during non-production")
	}
}
func updateDatabaseURL(config *app.Configuration) {
	if config.DatabaseUrl = os.Getenv("DATABASE_URL"); config.DatabaseUrl == "" {
		fmt.Println("DB URL: ", os.Getenv("DATABASE_URL"))
		flag.StringVar(&config.DatabaseUrl, "databaseUrl", "DEFAULT", "Database URL")
	}
}
func updateGenerateTestData(config *app.Configuration) {
	if config.GenerateTestData = strings.ToLower(os.Getenv("GENERATE_TEST_DATA")) == "true"; config.GenerateTestData == false {
		flag.BoolVar(&config.GenerateTestData, "generateTestData", true, "Generate test data")
	}
}
func updateDeleteLocalDatabase(config *app.Configuration) {
	if config.DeleteLocalDatabase = strings.ToLower(os.Getenv("DELETE_LOCAL_DATABASE")) == "true"; config.DeleteLocalDatabase == false {
		flag.BoolVar(&config.DeleteLocalDatabase, "deleteLocalDatabase", true, "Delete Local Database upon start")
	}
}
func updateAPIKey(config *app.Configuration) {
	if config.ApiKey = os.Getenv("API_KEY"); config.Port == "" {
		flag.StringVar(&config.ApiKey, "apiKey", "", "API key to use, will default to an empty string which will always deny")
	}
}

func localDBFileMaintenance(config *app.Configuration) error {
	if config.Production == false && config.DeleteLocalDatabase {
		localDatabaseName := "./" + config.LocalDB
		if f, err := os.Stat(localDatabaseName); errors.Is(err, os.ErrNotExist) {
			fmt.Println("File does not exist")
			return err
		} else {
			fmt.Println("File exists")
			err = os.Remove(localDatabaseName)
			if err != nil {
				return err
			}
			_, err := os.Create(localDatabaseName)
			if err != nil {
				return err
			}
			fmt.Printf("%s -> %d", f.Name(), f.Size())
		}
	}
	return nil
}

func connectToDatabase(config *app.Configuration) error {
	var db *gorm.DB
	var err error

	if config.Production {
		fmt.Println("Connecting to production database")
		db, err = gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{PrepareStmt: true})
	} else {
		fmt.Println("Connecting to non-production database")
		db, err = gorm.Open(sqlite.Open(config.LocalDB), &gorm.Config{PrepareStmt: true})
	}

	if err != nil {
		return errors.New("failed to connect database")
	}

	config.DB = db

	return err
}
