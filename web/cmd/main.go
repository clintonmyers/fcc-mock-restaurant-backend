package main

import (
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strings"
)

var maxIdle int
var maxOpenConn int
var lifetimeMinutes int
var production bool
var port string
var localDB string

func main() {
	flag.IntVar(&maxIdle, "maxIdle", 5, "Set max idle connections for database")
	flag.IntVar(&maxOpenConn, "maxOpenConn", 10, "Set max open connections for database")
	flag.IntVar(&lifetimeMinutes, "lifetimeMinutes", 30, "Set connection max lifetime")
	flag.BoolVar(&production, "production", false, "Is this a production run?")
	flag.StringVar(&port, "port", ":3000", "Port to use")
	flag.StringVar(&localDB, "localDB", "test.db", "Local database file used only during non-production")

	flag.Parse()

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	db, err := gorm.Open(sqlite.Open(localDB), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	configureDatabase(db)
	testing(db)
	api.Configuration = &models.Configuration{
		DB:         db,
		Production: production,
	}

	app := fiber.New()
	configureMiddleware(app)
	api.GetRoutes(app)

	app.Get("/", func(ctx *fiber.Ctx) error {
		var products []Product
		result := db.Find(&products)
		fmt.Printf("Num: %d\n", result.RowsAffected)
		fmt.Printf("error: %#v\n", result.Error)

		return ctx.Status(200).JSON(products)
	})

	log.Fatal(app.Listen(port))
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
