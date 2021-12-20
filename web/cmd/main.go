package main

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	dbHelper "github.com/clintonmyers/fcc-mock-restaurant-backend/db/helpers"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"log"
)

var globalConfig app.Configuration

func main() {

	loadConfiguration(&globalConfig)

	setDatabaseParameters(&globalConfig)
	configureMiddleware(globalConfig.WebApp)
	api.GetRoutes(globalConfig.WebApp, &globalConfig)

	err := dbHelper.MigrateDb(&globalConfig)
	if err != nil {
		log.Fatal(err)
	}

	if globalConfig.Production == false {
		// Will generate configuration data
		fmt.Println("Running testingGorm()")
		dbHelper.LoadTestData(&globalConfig)
	}

	log.Fatal(globalConfig.WebApp.Listen(globalConfig.Port))

}

// testing
