package main

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	dbHelper "github.com/clintonmyers/fcc-mock-restaurant-backend/db/helpers"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"log"
)

var globalConfig app.Configuration

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	loadConfiguration(&globalConfig)

	setDatabaseParameters(&globalConfig)
	configureMiddleware(globalConfig.WebApp)

	api.Repo = dbHelper.MainRepository{DB: globalConfig.DB}
	api.SetupRoutes(globalConfig.WebApp, &globalConfig)

	FatalIfErr(dbHelper.MigrateDB(&globalConfig))

	if globalConfig.GenerateTestData {
		fmt.Println("Loading Test Data")
		dbHelper.LoadTestData(&globalConfig)
	}

	//// Still working on graceful shutdown.
	//// I have tried this underneath Windows and have found it runs properly.
	//// Using GoLand and WSL this doesn't work properly
	//
	//c := make(chan os.Signal, 1) // Create signal channel
	//signal.Notify(c, os.Interrupt)
	//
	//go func() {
	//	_ = <-c
	//	fmt.Println("Gracefully shutdown starting")
	//	_ = globalConfig.WebApp.Shutdown()
	//}()
	
	FatalIfErr(globalConfig.WebApp.Listen(globalConfig.Port))

}

// testing
