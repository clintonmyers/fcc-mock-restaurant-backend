package main

import (
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
)

var globalConfig models.Configuration

func loadConfiguration(config *models.Configuration) {
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

	flag.Parse()

	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}
	fmt.Printf("ENV production: %s\n", os.Getenv("production"))

	// Setup DB connection
	if config.Production {
		log.Println("Connecting to production database")
		db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	} else {
		log.Println("Connecting to non-production database")
		db, err := gorm.Open(sqlite.Open(config.LocalDB), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	}

	//if !config.Production{
	//	fmt.Printf("ENV full configuration: %#v\n", config)
	//}

}

func main() {

	//globalConfig := models.Configuration{}
	loadConfiguration(&globalConfig)

	//globalConfig.DB = db

	configureDatabase(globalConfig.DB, &globalConfig)

	testing(globalConfig.DB)
	api.Configuration = &globalConfig

	app := fiber.New()
	configureMiddleware(app)
	api.GetRoutes(app)

	app.Get("/", func(ctx *fiber.Ctx) error {
		var products []Product
		result := globalConfig.DB.Find(&products)
		fmt.Printf("Num: %d\n", result.RowsAffected)
		fmt.Printf("error: %#v\n", result.Error)

		return ctx.Status(200).JSON(products)
	})

	log.Fatal(app.Listen(globalConfig.Port))
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func testing(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1) // find product with integer primary key
	fmt.Printf("%#v\n", product)

	db.First(&product, "code = ?", "D42") // find product with code D42
	fmt.Printf("%#v\n", product)

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	//db.Delete(&product, 1)
}
