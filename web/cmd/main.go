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

	err := dbHelper.MigrateDB(&globalConfig)
	if err != nil {
		log.Fatal(err)
	}

	if globalConfig.Production == false {
		// Will generate configuration data
		fmt.Println("Loading Test Data")
		dbHelper.LoadTestData(&globalConfig)
	}
	// Still working on graceful shutdown.
	// I believe this is a problem with WSL not properly registering OS signals
	//c := make(chan os.Signal, 1) // Create signal channel
	//signal.Notify(c, os.Interrupt)
	//
	//go func() {
	//	_ = <-c
	//	fmt.Println("Gracefully shutdown starting")
	//	_ = globalConfig.WebApp.Shutdown()
	//}()

	log.Fatal(globalConfig.WebApp.Listen(globalConfig.Port))

}

// testing
