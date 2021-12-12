package main

import (
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"gorm.io/gorm"
	"log"
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

func configureDatabase(db *gorm.DB, config *app.Configuration) {
	// Gets the underlying DB connection
	sqlDb, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	// configure the connection limits
	sqlDb.SetMaxIdleConns(config.MaxIdle)
	sqlDb.SetMaxOpenConns(config.MaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(config.LifetimeMinutes))
}
