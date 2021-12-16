package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func sayHello(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
	return c.SendString(msg)
}

func reviewsGrouping(app *fiber.App, config *app.Configuration) {
	g := app.Group("api", func(c *fiber.Ctx) error {
		m := c.Method()
		// Only allow idempotent methods without access token
		if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
			return c.Next()
		}

		// This needs to be first, so we will prevent an empty string from allowing access by default
		authToken := c.Get("AUTH_TOKEN", "")

		if authToken == "" || authToken != config.API_KEY {
			// We're returning a 404 because we want to avoid people scanning for apis that are guarded
			return c.Status(fiber.StatusNotFound).SendString("Cannot Find Requested Page")
		}

		return c.Next()
	})

	g.All("/testing/:name", sayHello)

	// TODO the 'reviews/' endpoints will need to have some kind of check to make sure users have access

	g.Get("/reviews/:id", func(c *fiber.Ctx) error {
		if id := c.Params("id"); id != "" {
			atoi, err := strconv.Atoi(id)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			var review models.Review
			config.DB.Find(&review, atoi)

			if uint(atoi) == review.ID {

				jsonReview := models.ReviewJson{}

				review.ToJsonVersion(&jsonReview)

				return c.Status(fiber.StatusOK).JSON(&jsonReview)
			}

		}
		return c.SendStatus(fiber.StatusNotFound)

	})
	g.Put("/reviews/:id", func(c *fiber.Ctx) error {
		// Can only update a review if we know which one we want to update
		id := c.Params("id")
		if id == "" || id == "0" {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		// TODO user verification
		// Parse what we're given by the user
		var json models.ReviewJson
		if err := c.BodyParser(&json); err == nil {

			if s, err := strconv.Atoi(id); err == nil {
				// Convert and save the update
				json.ID = uint(s)
				var review models.Review
				json.ToRegularVersion(&review)
				if config.DB.Save(&review).Error == nil {
					return c.Status(fiber.StatusCreated).SendString(c.Path())
				}
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})

	g.Post("/reviews", func(c *fiber.Ctx) error {
		// TODO user verification
		if c.Params("id", "0") == "0" {

			var json models.ReviewJson

			if err := c.BodyParser(&json); err == nil {
				// already looking at this in the params section
				//if json.ID != 0 {
				//	log.Println("Post for reviews/ api contains ID so not saving")
				//	return c.SendStatus(fiber.StatusBadRequest)
				//}

				var review models.Review
				json.ToRegularVersion(&review)
				config.DB.Create(&review)
				return c.Status(fiber.StatusCreated).SendString(strings.TrimSuffix(c.Route().Path, ":id") + "/" + strconv.Itoa(int(review.ID)))
			}
		}

		return c.SendStatus(fiber.StatusBadRequest)
	})

}

func GetRoutes(app *fiber.App, config *app.Configuration) {

	app.Get("/hello/:name", sayHello)
	reviewsGrouping(app, config)

}
