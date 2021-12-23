package helpers

import (
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm"
)

type MainRepository struct {
	DB *gorm.DB
}

func MigrateDB(config *app.Configuration) error {
	db := config.DB
	var err error

	err = db.AutoMigrate(&models.Company{})
	err = db.AutoMigrate(&models.Restaurant{})
	err = db.AutoMigrate(&models.User{})
	err = db.AutoMigrate(&models.UserRole{})
	err = db.AutoMigrate(&models.Testimonial{})
	err = db.AutoMigrate(&models.TestimonialImage{})
	err = db.AutoMigrate(&models.Address{})
	err = db.AutoMigrate(&models.MenuImage{})
	err = db.AutoMigrate(&models.Menu{})
	err = db.AutoMigrate(&models.Role{})
	err = db.AutoMigrate(&models.MenuItem{})

	return err
}

func LoadTestData(config *app.Configuration) error {
	db := config.DB

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

	//// Create a new company
	//
	//company := models.Company{
	//	Name:        "Fake Company 1",
	//	Description: "This is a fake company to test getting testimonials",
	//	Restaurants: nil,
	//}

	// Create a Menu

	italianMenu := models.Menu{
		Name:      "Menu1",
		Active:    true,
		MenuItems: nil,
	}

	// Create MenuItems

	pizzaItem := models.MenuItem{
		Name:        "Pizza",
		Price:       1000,
		Description: "Pepperoni Pizza",
		Images: []models.MenuImage{
			{
				ImageURL: "https://giphy.com/gifs/pizzahut-happy-cute-3oz8xH46dD1DSx3vNK",
			},
			{
				ImageURL: "https://giphy.com/gifs/pizza-food-galaxy-GFKZLvfBiyrN6",
			},
		},
	}
	pastaItem := models.MenuItem{
		Name:        "Spaghetti",
		Price:       500,
		Description: "Plate of spaghetti",
		Images: []models.MenuImage{
			{
				ImageURL: "https://media.giphy.com/media/oFy0DysfL2nGG0W6Gr/giphy.gif",
			},
		},
		MenuID: 0,
	}
	italianMenu.MenuItems = []models.MenuItem{pizzaItem, pastaItem}

	americanMenu := models.Menu{
		Name:   "American Classic Cuisine",
		Active: true,
		MenuItems: []models.MenuItem{
			models.MenuItem{
				Name:        "Chicken and Waffles",
				Price:       1500,
				Description: "Classic Chicken and Waffles dish",
				Images: []models.MenuImage{
					{
						ImageURL: "https://media.giphy.com/media/WChsdz6hvhiXC/giphy.gif",
					},
				},
			},
			{
				//Model:       gorm.Model{},
				Name:        "BBQ Ribs",
				Price:       1500,
				Description: "Classic American smoked ribs",
				Images: []models.MenuImage{
					models.MenuImage{
						ImageURL: "https://media.giphy.com/media/YBZ0UWAaqWes5PrL8x/giphy.gif",
					},
				},
				MenuID: 0,
			},
		},
	}
	// Create a Restaurant
	rest1 := models.Restaurant{
		Addresses: []models.Address{addy},
		Menus:     []models.Menu{italianMenu},
	}

	rest2 := models.Restaurant{
		Addresses: []models.Address{addy2},
		Menus:     []models.Menu{americanMenu},
	}

	imgUrl := models.TestimonialImage{
		URL: "https://giphy.com/gifs/rick-roll-gotcha-mod-miny-kFgzrTt798d2w",
	}

	imgUrl2 := models.TestimonialImage{
		URL: "https://giphy.com/gifs/officialblueytv-yes-wow-bluey-fXiEx2PQF1bBXYi7bG",
	}

	imgUrl3 := models.TestimonialImage{
		URL: "https://giphy.com/gifs/officialblueytv-yes-wow-bluey-StdOGnXLUnxSC3Hx40",
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

	db.FullSaveAssociations = true
	company := models.Company{
		Name:        "Fake Company 1",
		Description: "First Fake Company",
		Restaurants: []models.Restaurant{rest1, rest2},
	}
	//company.Restaurants = []models.Restaurant{rest1, rest2}
	//db.Debug().Clauses(clause.OnConflict{
	//	DoNothing: true,
	//}).Save(&company)
	//db.Save()
	db.Create(&company)
	if db.Error != nil {
		return db.Error
	}
	return nil

}

func FullCompany() models.Company {
	return models.Company{
		Name:        "TEST1",
		Description: "Test 1",
		Restaurants: []models.Restaurant{
			{
				Addresses: []models.Address{
					{
						Address: "123 Fake St",
						City:    "Dallas",
						State:   "TX",
						Zip:     "12345",
						Active:  true,
					},
				},
				Testimonials: []models.Testimonial{
					{
						Header: "header",
						Body:   "body",
						ImageUrls: []models.TestimonialImage{
							{
								URL: "LINK_1",
							},
						},
						RestaurantID: 0,
					},
				},
				Menus: []models.Menu{
					{
						Name:   "Menu1",
						Active: true,
						MenuItems: []models.MenuItem{{
							Name:        "MenuItemName",
							Price:       12345,
							Description: "Description",
							Images: []models.MenuImage{
								{
									ImageURL: "LINK_2",
								},
							},
						}},
					},
				},
			},
		},
	}
}
