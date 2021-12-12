package main

import (
	"flag"
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
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

var globalConfig app.Configuration

func loadConfiguration(config *app.Configuration) {
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

	if config.GenerateTestData = strings.ToLower(os.Getenv("GENERATE_TEST_DATA")) == "true"; config.GenerateTestData == false {
		flag.BoolVar(&config.GenerateTestData, "generateTestData", true, "Generate test data")
	}

	flag.Parse()

	if !strings.HasPrefix(config.Port, ":") {
		config.Port = ":" + config.Port
	}

	// Setup DB connection
	if config.Production {
		fmt.Println("Connecting to production database")
		db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	} else {
		fmt.Println("Connecting to non-production database")
		db, err := gorm.Open(sqlite.Open(config.LocalDB), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			panic("failed to connect database")
		}
		config.DB = db
	}

	// Setup the web app
	config.WebApp = fiber.New()

	fmt.Printf("ENV production: %s\n", os.Getenv("production"))

	if !config.Production {
		fmt.Printf("ENV full configuration: %#v\n", config)
	}

}

func main() {

	//globalConfig := models.Configuration{}
	loadConfiguration(&globalConfig)

	//globalConfig.DB = db

	configureDatabase(globalConfig.DB, &globalConfig)

	testing(globalConfig.DB)
	api.Configuration = &globalConfig
	app := globalConfig.WebApp
	//app := fiber.New()
	configureMiddleware(app)
	api.GetRoutes(app)
	testingGorm(&globalConfig)

	//app.Get("/", func(ctx *fiber.Ctx) error {
	//	var products []Product
	//	result := globalConfig.DB.Find(&products)
	//	fmt.Printf("Num: %d\n", result.RowsAffected)
	//	fmt.Printf("error: %#v\n", result.Error)
	//
	//	return ctx.Status(200).JSON(products)
	//})

	log.Fatal(app.Listen(globalConfig.Port))
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func testingGorm(c *app.Configuration) {
	db := c.DB

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Restaurant{})
	db.AutoMigrate(&models.MenuItem{})
	db.AutoMigrate(&models.Review{
		Model:      gorm.Model{},
		Restaurant: nil,
		Header:     "",
		Body:       "",
		Images:     make([]string, 0),
	})

	r := models.Restaurant{
		Address: models.Address{
			Address: "123 fake street",
			City:    "Fakeville",
			State:   "TX",
			Zip:     "12345",
		},
		Users: make([]*models.User, 0, 10),
	}

	u := models.User{
		Role:      0,
		Username:  "User1",
		FirstName: "first",
		LastName:  "last",
		Address: models.Address{
			Address: "456 Fake Street",
			City:    "Fakeville",
			State:   "Tx",
			Zip:     "78901",
		},
	}
	u2 := models.User{
		Role:      0,
		Username:  "User2",
		FirstName: "first2",
		LastName:  "last2",
		Address: models.Address{
			Address: "4561 Fake Street",
			City:    "Fakeville",
			State:   "Tx",
			Zip:     "78901",
		},
	}

	db.Save(&u2)
	r.Users = append(r.Users, &u)

	menuItem := models.MenuItem{
		Name:        "Pizza",
		Price:       3.50,
		Description: "Pepperoni Pizza",
		Image:       "LINK_TO_PIZZA",
		Restaurant:  make([]*models.Restaurant, 0, 10),
	}
	menuItem.Restaurant = append(menuItem.Restaurant, &r)

	review := models.Review{
		Restaurant: &r,
		Header:     "Review Header #1",
		Body:       "This is the main content of the review. It's got a lot of text inside of it",
		Images:     make([]string, 0, 1),
	}

	review.Images = append(review.Images, "LINK_TO_REVIEW_IMAGE")

	// User test
	rowsAffected := db.Save(&u).RowsAffected
	fmt.Println("Table users rows affected: ", rowsAffected)

	// Restaurant test
	rowsAffected = db.Save(&r).RowsAffected
	fmt.Println("Table restaurants rows affected: ", rowsAffected)

	// Menu Item
	rowsAffected = db.Save(&menuItem).RowsAffected
	fmt.Println("Table MenuItem rows affected: ", rowsAffected)

	// Reviews
	rowsAffected = db.Save(&review).RowsAffected
	fmt.Println("Table Review rows affected: ", rowsAffected)

	app := c.WebApp
	app.Get("/", func(c *fiber.Ctx) error {
		//var user models.User
		//db.First(&user)
		////db.Find(&user, )

		//var rr models.Restaurant
		//db.First(&rr)
		//db.Find()
		//var restaurant models.Restaurant
		//db.Find(&restaurant)

		var rr models.Restaurant
		//db.First(&rr)
		//db.Model(&rr).Preload("Users")
		db.Preload("Users").First(&rr)
		//db.Where("users in ?", []uint{1}).First(&rr)
		//fmt.Printf("%#v\n", rr)

		return c.Status(fiber.StatusOK).JSON(&rr)
	})

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
