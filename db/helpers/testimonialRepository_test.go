package helpers

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
	"testing"
)

func TestMainRepository_SaveSanityCheck(t *testing.T) {
	// Setup
	tempDir := t.TempDir()

	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	//var mainRepo TestimonialRepository
	//mainRepo = &MainRepository{DB: db}
	//_ = mainRepo

	LoadTestData(&app.Configuration{DB: db})

	//
	//db.Save(&company1)
	//
	//if db.Error != nil {
	//	t.Fail()
	//}

	/*
		Basic logic to make sure that our test data is correct
	*/

	LoadTestData(&app.Configuration{DB: db})

	var testCompany models.Company
	// The following work but don't pull in all of the associations for a complete data object
	//db.Find(&testCompany, 1)
	//db.Preload(clause.Associations).Preload("Menus").First(&testCompany)
	db.Preload(clause.Associations).Preload("Restaurants.Testimonials." + clause.Associations).Preload("Restaurants.Menus.MenuItems." + clause.Associations).First(&testCompany)

	// Basic sanity check to make sure that we can actually save all the data
	if testCompany.ID != 1 ||
		len(testCompany.Restaurants) != 2 ||
		len(testCompany.Restaurants[0].Testimonials) != 3 ||
		len(testCompany.Restaurants[1].Testimonials) != 2 ||
		len(testCompany.Restaurants[0].Menus) != 1 ||
		len(testCompany.Restaurants[1].Menus) != 1 {
		t.Fail()
	}

}

func TestMainRepository_SaveTestimonial(t *testing.T) {

	tempDir := t.TempDir()
	// Still need to make sure that we close the file resource or we'll have a problem with the directory being used
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	var mainRepo TestimonialRepository
	mainRepo = &MainRepository{DB: db}

	//// Migrate the DB tables
	//db.AutoMigrate(&models.Testimonial{})
	//db.AutoMigrate(&models.TestimonialImage{})
	//db.AutoMigrate(&models.Restaurant{})
	//db.AutoMigrate(&models.Address{})
	_ = mainRepo

	// -------------------------------- //
	// -------------------------------- //
	// Create the data
	// -------------------------------- //
	// -------------------------------- //

	// This is testing the restaurant more than the testimonial
	//// Create an Address
	//addy := models.Address{
	//	Address: "123 Fake Street",
	//	City:    "Dallas",
	//	State:   "TX",
	//	Zip:     "12345",
	//	Active:  false,
	//}
	//
	//// Create a Restaurant
	//
	//rest := models.Restaurant{
	//	Addresses:    []models.Address{addy},
	//	CompanyId:    0,
	//	Testimonials: []models.Testimonial{},
	//}

	imgUrl := models.TestimonialImage{
		URL: "TEST_LINK",
	}
	// Create a Testimonial
	testimonial := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl},
	}
	// This is testing the Restaurant API not strictly the testimonial one
	//rest.Testimonials = append(rest.Testimonials, testimonial)

	// -------------------------------- //
	// -------------------------------- //
	// Save to the database
	// -------------------------------- //
	// -------------------------------- //

	//count := db.Save(&rest).RowsAffected
	//if db.Error != nil {
	//	t.Fatal(db.Error)
	//}
	//fmt.Println("resty count: ", count)

	c, err := mainRepo.SaveTestimonial(&testimonial)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("COUNT: ", c)

	//var restaurant models.Restaurant
	//db.First(&restaurant)
	//fmt.Println(restaurant)

	// -------------------------------- //
	// -------------------------------- //
	// Make the actual testing happen
	// -------------------------------- //
	// -------------------------------- //

	var testy models.Testimonial
	//db.Preload(clause.Associations).First(&testy)
	mainRepo.GetTestimonialById(&testy, testimonial.ID)
	fmt.Println(testy)

	if testy.ID != 1 {
		fmt.Println("This should have been saved properly")
		fmt.Printf("OUTPUT: %#v", testy)
		t.Fail()
	}
	if len(testy.ImageUrls) < 1 || testy.ImageUrls[0].ID != 1 {
		fmt.Println("should have had some image urls")
		t.Fail()
	}
	//fmt.Printf("\n %#v\n", testy.ImageUrls)

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //
	// -------------------------------- //

}
