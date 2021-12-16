package main

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"log"
)

var globalConfig app.Configuration

func main() {

	loadConfiguration(&globalConfig)

	setDatabaseParameters(&globalConfig)
	configureMiddleware(globalConfig.WebApp)
	api.GetRoutes(globalConfig.WebApp, &globalConfig)

	err := migrateDb(&globalConfig)
	if err != nil {
		log.Fatal(err)
	}

	if globalConfig.Production == false {
		// Will generate configuration data
		fmt.Println("Running testingGorm()")
		testingGorm(&globalConfig)
	}

	log.Fatal(globalConfig.WebApp.Listen(globalConfig.Port))

}
func migrateDb(c *app.Configuration) error {
	if c.AutoMigrate == false {
		fmt.Println("Not running auto migration")
		return nil
	}
	fmt.Println("Running auto migration")

	var err error
	db := c.DB

	err = db.AutoMigrate(&models.Address{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Company{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.MenuItem{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Restaurant{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Review{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.User{})

	return err
}

// testing
func testingGorm(c *app.Configuration) {
	db := c.DB

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
	test := make([]string, 0)
	test = append(test, "A")
	test = append(test, "b")
	review := models.Review{
		Restaurant: &r,
		Header:     "Review Header #1",
		Body:       "This is the main content of the review. It's got a lot of text inside of it",
		Images:     "",
	}
	review.AppendString("LINK_TO_REVIEW_IMAGE")
	fmt.Println("review.GetStringArr(): ", review.GetStringArr())
	review.AppendString("LINK_TO_REVIEW_IMAGE2")
	//review.Images = append(review.Images, "LINK_TO_REVIEW_IMAGE")
	//review.Images = append(review.Images, "LINK_TO_REVIEW_IMAGE2")
	fmt.Printf("The Review: %v\n", review)
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
	db.Save(&review)
	fmt.Println("Table Review rows affected: ", rowsAffected)

	//app := c.WebApp
	//app.Get("/", func(c *fiber.Ctx) error {
	//
	//	var rr models.Restaurant
	//	db.Preload("Users").First(&rr)
	//
	//	return c.Status(fiber.StatusOK).JSON(&rr)
	//})

}
