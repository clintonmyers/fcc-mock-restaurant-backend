package main

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

var globalConfig app.Configuration

func main() {

	//globalConfig := models.Configuration{}
	loadConfiguration(&globalConfig)

	//globalConfig.DB = db

	configureDatabase(&globalConfig)

	testing(globalConfig.DB)
	api.Configuration = &globalConfig
	app := globalConfig.WebApp
	//app := fiber.New()
	configureMiddleware(app)
	api.GetRoutes(app)
	testingGorm(&globalConfig)

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
