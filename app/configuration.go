package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type Configuration struct {
	DB     *gorm.DB
	WebApp *fiber.App
	Store  *session.Store

	/*
		DATABASE CONFIGURATION
	*/
	MaxIdle     int
	MaxOpenConn int
	LocalDB     string
	DatabaseUrl string
	AutoMigrate bool

	LifetimeMinutes int
	// Web server, this app, port
	Port string

	// Will choose whether to connect to remote or test database
	Production bool
	// Should we generate Test Data?
	GenerateTestData bool
	// Will only happen if production is false as well
	DeleteLocalDatabase bool

	/*
		AUTHENTICATION
	*/
	// No longer used, but was an api key for authentication
	ApiKey string

	// Google authentication keys
	GoogleOAuthKey    string
	GoogleOAuthSecret string
	// Where we ask the oauth provider to call us back at, this app
	CallbackUrl string
	// If true will accept a form with 'user' & 'pass' options being passed in
	SimulateOAuth bool
	// The username for simulated oauth user
	SimulatedUser string
	// the password for simulated oauth user
	SimulatedPassword string
	// The secret key that will be used to encrypt the jwt token
	OAuthSecret string
	// Where do we redirect a user after they have successfully authenticated
	AuthRedirect string

	/*
		SESSION LOCATION, will still be used for authentication redirect cookie
	*/
	// Where do we store the session
	SessionLocation string
	SessionName     string
}

const (
	API_KEY_OS      string = "API_KEY"
	API_KEY_FLAG    string = "apiKey"
	API_KEY_DEFAULT string = ""
	API_KEY_HEADER  string = "API_KEY"

	DELETE_LOCAL_DB_OS      string = "DELETE_LOCAL_DATABASE"
	DELETE_LOCAL_DB_FLAG    string = "deleteLocalDatabase"
	DELETE_LOCAL_DB_DEFAULT bool   = true

	GENERATE_TEST_DATA_OS      string = "GENERATE_TEST_DATA"
	GENERATE_TEST_DATA_FLAG    string = "generateTestData"
	GENERATE_TEST_DATA_DEFAULT bool   = false

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
	PORT_DEFAULT string = ":8088"

	GOOGLE_AUTH_KEY_OS   string = "GOOGLE_AUTH_KEY"
	GOOGLE_AUTH_KEY_FLAG string = "googleAuthKey"
	GOOGLE_AUTH_DEFAULT  string = ""

	GOOGLE_SECRET_KEY_OS   string = "GOOGLE_SECRET_KEY"
	GOOGLE_SECRET_KEY_FLAG string = "googleSecretKey"
	GOOGLE_SECRET_DEFAULT  string = ""

	CALLBACK_URL_OS   string = "REDIRECT_URL"
	CALLBACK_URL_FLAG string = "redirectUrl"
	CALLBACK_DEFAULT  string = "http://127.0.0.1:8088/auth/google/callback"

	SIMULATE_OAUTH_OS      string = "SIMULATE_OAUTH"
	SIMULATE_OAUTH_FLAG    string = "simulateOauth"
	SIMULATE_OAUTH_DEFAULT bool   = false

	SIMULATED_USER_OS      string = "SIMULATED_USER"
	SIMULATED_USER_FLAG    string = "simulatedUser"
	SIMULATED_USER_DEFAULT string = ""

	SIMULATED_PASSWORD_OS      string = "SIMULATED_PASSWORD"
	SIMULATED_PASSWORD_FLAG    string = "simulatedPassword"
	SIMULATED_PASSWORD_DEFAULT string = ""

	OAUTH_SECRET_OS      string = "OAUTH_SECRET"
	OAUTH_SECRET_FLAG    string = "oauthSecret"
	OAUTH_SECRET_DEFAULT string = ""

	SESSION_LOCATION_OS      string = "SESSION_LOCATION"
	SESSION_LOCATION_FLAG    string = "sessionLocation"
	SESSION_LOCATION_DEFAULT string = "cookie"

	SESSION_NAME_OS      string = "SESSION_NAME"
	SESSION_NAME_FLAG    string = "sessionName"
	SESSION_NAME_DEFAULT string = "session_id"

	AUTH_REDIRECT_OS      string = "AUTH_REDIRECT"
	AUTH_REDIRECT_FLAG    string = "authRedirect"
	AUTH_REDIRECT_DEFAULT string = "http://localhost:3000/login"
)
