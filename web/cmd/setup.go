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
	//updateGenerateLocalData(config)
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
	if config.Port = os.Getenv(app.PORT_OS); config.Port == "" {
		flag.StringVar(&config.Port, app.PORT_FLAG, app.PORT_DEFAULT, "Port to use")
	}
}
func updateAutoMigrate(config *app.Configuration) {
	if config.AutoMigrate = strings.ToLower(os.Getenv(app.AUTO_MIGRATE_OS)) == "true"; config.AutoMigrate == false {
		flag.BoolVar(&config.AutoMigrate, app.AUTO_MIGRATE_FLAG, app.AUTO_MIGRATE_DEFAULT, "Should we auto migrate the database?")
	}
}

func updateProductionSetting(config *app.Configuration) {
	if config.Production = strings.ToLower(os.Getenv(app.PRODUCTION_OS)) == "true"; config.Production == false {
		flag.BoolVar(&config.Production, app.PRODUCTION_FLAG, app.PRODUCTION_DEFAULT, "Is this a production run?")
	}
}

//func updateGenerateLocalData(config *app.Configuration) {
//	if config.GenerateLocalData = strings.ToLower(os.Getenv("GENERATE_LOCAL_DATA")) == "true"; config.GenerateLocalData == false {
//		flag.BoolVar(&config.GenerateLocalData, "generateData", false, "Do we generate local data?")
//	}
//}

func updateMaxIdle(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv(app.MAX_IDLE_OS)); err == nil {
		config.MaxIdle = s
	} else {
		flag.IntVar(&config.MaxIdle, app.MAX_IDLE_FLAG, app.MAX_IDLE_DEFAULT, "Set max idle connections for database")
	}
}
func updateMaxOpenConnections(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv(app.MAX_OPEN_CONN_OS)); err == nil {
		config.MaxOpenConn = s
	} else {
		flag.IntVar(&config.MaxOpenConn, app.MAX_OPEN_CONN_FLAG, app.MAX_OPEN_CONN_DEFAULT, "Set max open connections for database")
	}
}
func updateLifetimeMinutes(config *app.Configuration) {
	if s, err := strconv.Atoi(os.Getenv(app.LIFETIME_MINUTES_OS)); err == nil {
		config.LifetimeMinutes = s
	} else {
		flag.IntVar(&config.LifetimeMinutes, app.LIFETIME_MINUTES_FLAG, app.LIFETIME_MINUTES_DEFAULT, "Set max open connections for database")
	}
}
func updateLocalDBSetting(config *app.Configuration) {
	if config.LocalDB = os.Getenv(app.LOCAL_DB_OS); config.LocalDB == "" {
		flag.StringVar(&config.LocalDB, app.LOCAL_DB_FLAG, app.LOCAL_DB_DEFAULT, "Local database file used only during non-production")
	}
}
func updateDatabaseURL(config *app.Configuration) {
	if config.DatabaseUrl = os.Getenv("DATABASE_URL"); config.DatabaseUrl == "" {
		fmt.Println("DB URL: ", os.Getenv(app.DATABASE_URL_OS))
		flag.StringVar(&config.DatabaseUrl, app.DATABASE_URL_FLAG, app.DATABASE_URL_DEFAULT, "Database URL")
	}
}
func updateGenerateTestData(config *app.Configuration) {
	if config.GenerateTestData = strings.ToLower(os.Getenv(app.GENERATE_TEST_DATA_OS)) == "true"; config.GenerateTestData == false {
		flag.BoolVar(&config.GenerateTestData, app.GENERATE_TEST_DATA_FLAG, app.GENERATE_TEST_DATA_DEFAULT, "Generate test data")
	}
}
func updateDeleteLocalDatabase(config *app.Configuration) {
	if config.DeleteLocalDatabase = strings.ToLower(os.Getenv(app.DELETE_LOCAL_DB_OS)) == "true"; config.DeleteLocalDatabase == false {
		flag.BoolVar(&config.DeleteLocalDatabase, app.DELETE_LOCAL_DB_FLAG, app.DELETE_LOCAL_DB_DEFAULT, "Delete Local Database upon start")
	}
}
func updateAPIKey(config *app.Configuration) {
	if config.ApiKey = os.Getenv(app.API_KEY_OS); config.ApiKey == "" {
		flag.StringVar(&config.ApiKey, app.API_KEY_FLAG, app.API_KEY_DEFAULT, "API key to use, will default to an empty string which will always deny")
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
