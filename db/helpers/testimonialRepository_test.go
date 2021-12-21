package helpers

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"testing"
)

func TestMainRepository_SaveSanityCheck(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	// Still need to make sure that we close the file resource or we'll have a problem with the directory being used
	file, db, err := CreateAndMigrateTempDB(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	var mainRepo TestimonialRepository
	mainRepo = &MainRepository{DB: db}
	_ = mainRepo

	// -------------------------------- //
	// -------------------------------- //
	// Create the data
	// -------------------------------- //
	// -------------------------------- //

	// Create the addresses
	addy := models.Address{
		Address: "123 Fake Street",
		City:    "Dallas",
		State:   "TX",
		Zip:     "12345",
		Active:  true,
	}
	addy2 := models.Address{
		Address: "321 Fake Street",
		City:    "Dallas",
		State:   "TX",
		Zip:     "54321",
		Active:  true,
	}

	//addy3 := models.Address{
	//	Address: "456 Fake Street",
	//	City:    "Dallas",
	//	State:   "TX",
	//	Zip:     "12345",
	//	Active:  true,
	//}

	// Create a new company

	company := models.Company{
		Name:        "Fake Company 1",
		Description: "This is a fake company to test getting testimonials",
		Restaurants: nil,
	}
	// Create a Restaurant

	rest1 := models.Restaurant{
		Addresses: []models.Address{addy},
		//Testimonials: []models.Testimonial{},
	}

	rest2 := models.Restaurant{
		Addresses: []models.Address{addy2},
		//Testimonials: nil,
	}

	// Assign the restaurants to the company
	//company.Restaurants = []models.Restaurant{rest1, rest2}

	imgUrl := models.TestimonialImage{
		URL: "TEST_LINK_1",
	}

	imgUrl2 := models.TestimonialImage{
		URL: "TEST_LINK_2",
	}

	imgUrl3 := models.TestimonialImage{
		URL: "TEST_LINK_3",
	}
	// Create a Testimonial
	testimonial := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl, imgUrl2, imgUrl3},
	}
	testimonial2 := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl, imgUrl2, imgUrl3},
	}
	testimonial3 := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl, imgUrl2, imgUrl3},
	}
	testimonial4 := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl, imgUrl2, imgUrl3},
	}
	testimonial5 := models.Testimonial{
		//Restaurant: rest,
		Header:    "",
		Body:      "",
		ImageUrls: []models.TestimonialImage{imgUrl, imgUrl2, imgUrl3},
	}

	rest1.Testimonials = []models.Testimonial{testimonial, testimonial2, testimonial3}
	rest2.Testimonials = []models.Testimonial{testimonial4, testimonial5}

	//db.FullSaveAssociations = true
	db.Save(&rest1)
	db.Save(&rest2)
	company.Restaurants = []models.Restaurant{rest1, rest2}
	db.Save(&company)
	fmt.Printf("TESTING COMPANY: %#v\n", company)

	var company2 models.Company
	db.Find(&company2, 1)

	// Basic sanity check to make sure that we can actually save all the data
	if company2.ID != 1 &&
		len(company2.Restaurants) != 2 &&
		len(company2.Restaurants[0].Testimonials) != 3 &&
		len(company2.Restaurants[1].Testimonials) != 2 {
		t.Fail()
	}

	//db.Preload(clause.Associations).Find(&company2, 1)
	//db.Preload("Restaurants").Preload("Testimonials").Find(&company2, 1)
	//db.Preload("Restaurant.Testimonials").Find(&company2, 1)
	//db.Where("RestaurantID")
	//db.Preload("Testimonials").Preload("Restaurants").Find(&company2, 1)
	// This is testing the Restaurant API not strictly the testimonial one
	//rest.Testimonials = append(rest.Testimonials, testimonial)

	// -------------------------------- //
	// -------------------------------- //
	// Save to the database
	// -------------------------------- //
	// -------------------------------- //
	//
	//count := db.Save(&rest).RowsAffected
	//if db.Error != nil {
	//	t.Fatal(db.Error)
	//}
	//fmt.Println("resty count: ", count)
	//
	//c, err := mainRepo.SaveTestimonial(&testimonial)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//fmt.Println("COUNT: ", c)
	//
	//var restaurant models.Restaurant
	//db.First(&restaurant)
	//fmt.Println(restaurant)

	// -------------------------------- //
	// -------------------------------- //
	// Make the actual testing happen
	// -------------------------------- //
	// -------------------------------- //

	//var testy models.Testimonial
	////db.Preload(clause.Associations).First(&testy)
	//mainRepo.GetTestimonialById(&testy, testimonial.ID)
	//fmt.Println(testy)
	//
	//if testy.ID != 1 {
	//	fmt.Println("This should have been saved properly")
	//	fmt.Printf("OUTPUT: %#v", testy)
	//	t.Fail()
	//}
	//if len(testy.ImageUrls) < 1 || testy.ImageUrls[0].ID != 1 {
	//	fmt.Println("should have had some image urls")
	//	t.Fail()
	//}
	//fmt.Printf("\n %#v\n", testy.ImageUrls)

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
