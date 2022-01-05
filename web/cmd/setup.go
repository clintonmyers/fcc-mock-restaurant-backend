package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/session"
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
	updateMaxIdle(config)
	updateMaxOpenConnections(config)
	updateLifetimeMinutes(config)
	updateLocalDBSetting(config)
	updateDatabaseURL(config)
	updateGenerateTestData(config)
	updateDeleteLocalDatabase(config)
	updateGoogleOAuthSettings(config)
	updateSimulateOauth(config)
	updateOAuthSecret(config)
	updateSessionConfigs(config)

	// <TEMP CONFIGS>: these should be removed when we're actually doing real authentication
	updateAPIKey(config)
	// </TEMP CONFIGS>

	flag.Parse()

	createSessionStore(config)

	if err := validateConfiguration(config); err != nil {
		log.Fatal(err)
	}
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

func validateConfiguration(config *app.Configuration) error {
	if config.OAuthSecret == "" {
		return errors.New("Cannot have an empty secret")
	}

	if config.SimulateOAuth && (config.SimulatedUser == "" || config.SimulatedPassword == "") {
		return errors.New("Cannot have simulated OAuth annd empty username/password")
	}

	return nil
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

func updateGoogleOAuthSettings(config *app.Configuration) {
	// For Google OAuth we need to setup the keys and the secret as well as a callback URL for this to return to
	if config.GoogleOAuthKey = strings.ToLower(os.Getenv(app.GOOGLE_AUTH_KEY_OS)); config.GoogleOAuthKey == "" {
		flag.StringVar(&config.GoogleOAuthKey, app.GOOGLE_AUTH_KEY_FLAG, app.GOOGLE_AUTH_DEFAULT, "Used to set the Google OAuth key")
	}
	if config.GoogleOAuthSecret = strings.ToLower(os.Getenv(app.GOOGLE_SECRET_KEY_OS)); config.GoogleOAuthSecret == "" {
		flag.StringVar(&config.GoogleOAuthSecret, app.GOOGLE_SECRET_KEY_FLAG, app.GOOGLE_SECRET_DEFAULT, "Used to set the Google OAuth Secret")
	}
	if config.CallbackUrl = strings.ToLower(os.Getenv(app.CALLBACK_URL_OS)); config.CallbackUrl == "" {
		flag.StringVar(&config.CallbackUrl, app.CALLBACK_URL_FLAG, app.CALLBACK_DEFAULT, "Used to set the Google OAuth callback")
	}
}
func updateSimulateOauth(config *app.Configuration) {
	if config.SimulateOAuth = strings.ToLower(os.Getenv(app.SIMULATE_OAUTH_OS)) == "true"; config.DeleteLocalDatabase == false {
		flag.BoolVar(&config.SimulateOAuth, app.SIMULATE_OAUTH_FLAG, app.SIMULATE_OAUTH_DEFAULT, "Should we simulate Oauth with a /login/")
	}
	if config.SimulatedUser = strings.ToLower(os.Getenv(app.SIMULATED_USER_OS)); config.SimulatedUser == "" {
		flag.StringVar(&config.SimulatedUser, app.SIMULATED_USER_FLAG, app.SIMULATED_USER_DEFAULT, "Simulated User")
	}
	if config.SimulatedPassword = strings.ToLower(os.Getenv(app.SIMULATED_PASSWORD_OS)); config.SimulatedPassword == "" {
		flag.StringVar(&config.SimulatedPassword, app.SIMULATED_PASSWORD_FLAG, app.SIMULATED_PASSWORD_DEFAULT, "Simulated Password")
	}
}
func updateAPIKey(config *app.Configuration) {
	if config.ApiKey = os.Getenv(app.API_KEY_OS); config.ApiKey == "" {
		flag.StringVar(&config.ApiKey, app.API_KEY_FLAG, app.API_KEY_DEFAULT, "API key to use, will default to an empty string which will always deny")
	}
}

func updateOAuthSecret(config *app.Configuration) {
	if config.OAuthSecret = os.Getenv(app.OAUTH_SECRET_OS); config.OAuthSecret == "" {
		flag.StringVar(&config.OAuthSecret, app.OAUTH_SECRET_FLAG, app.OAUTH_SECRET_DEFAULT, "API key to use, will default to an empty string which will always deny")
	}
}

func updateSessionConfigs(config *app.Configuration) {
	if config.SessionLocation = os.Getenv(app.SESSION_LOCATION_OS); config.SessionLocation == "" {
		flag.StringVar(&config.SessionLocation, app.SESSION_LOCATION_FLAG, app.SESSION_LOCATION_DEFAULT, "session location [header/cookie]")
	}
	if config.SessionName = os.Getenv(app.SESSION_NAME_OS); config.SessionName == "" {
		flag.StringVar(&config.SessionName, app.SESSION_NAME_FLAG, app.SESSION_NAME_DEFAULT, "session name cookie:session_id")
	}
}

func createSessionStore(config *app.Configuration) {
	config.Store = session.New(session.Config{
		KeyLookup: fmt.Sprintf("%s:%s", config.SessionLocation, config.SessionName),
	})
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
