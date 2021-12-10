package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
)

var testimonials = []models.Testimonial{
	{
		ID:      0,
		Title:   "Testimonial 1",
		Comment: "Comment #1. There is so much to say about this once!!",
		Rating:  5,
	},
	{
		ID:      1,
		Title:   "Testimonial 2",
		Comment: "Comment #2. There is so much to say about this once!!",
		Rating:  5,
	},
	{
		ID:      2,
		Title:   "Testimonial 3",
		Comment: "Comment #3",
		Rating:  3,
	}, {
		ID:      3,
		Title:   "Testimonial 4",
		Comment: "Comment #4",
		Rating:  4,
	}, {
		ID:      4,
		Title:   "Testimonial 5",
		Comment: "Comment #5",
		Rating:  1,
	},
}

var Configuration *models.Configuration

func GetRoutes(app *fiber.App) {
	fmt.Println("HERE")
	Configuration.DB.AutoMigrate(&models.Testimonial{})

	if !Configuration.Production {
		for _, t := range testimonials {
			Configuration.DB.Create(&t)
		}

	}

	app.Get("/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	g := app.Group("api")

	g.Get("testimonials", func(c *fiber.Ctx) error {
		var ts []models.Testimonial
		result := Configuration.DB.Find(&ts)
		fmt.Printf("Num: %d\n", result.RowsAffected)
		fmt.Printf("error: %#v\n", result.Error)

		return c.Status(fiber.StatusOK).JSON(ts)
	})

}
