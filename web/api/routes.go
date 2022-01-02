package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/db/helpers"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

func sayHello(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
	return c.SendString(msg)
}

func getJwtFunction(config *app.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.OAuthSecret),
		Filter: func(ctx *fiber.Ctx) bool {
			m := ctx.Method()
			if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
				return true
			}
			return false
		},
	})
}

func apiKeyAuth(config *app.Configuration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		localTest := ctx.Locals("user")
		fmt.Println(localTest)
		m := ctx.Method()
		// Only allow idempotent methods without access token
		if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
			return ctx.Next()
		}

		// This needs to be first, so we will prevent an empty string from allowing access by default
		authToken := ctx.Get(app.API_KEY_HEADER, "")

		if authToken == "" || authToken != config.ApiKey {
			// We're returning a 404 because we want to avoid people scanning for apis that are guarded
			return ctx.Status(fiber.StatusNotFound).SendString("Cannot Find Requested Page")
		}

		return ctx.Next()
	}
}
func SetupRoutes(fiberApp *fiber.App, config *app.Configuration) {
	fmt.Println("Setting Routes")

	fiberApp.Get("/hello/:name", sayHello)

	if config.SimulateOAuth {
		fiberApp.Post("/login", simulatedLogin(config))
	}

	api := fiberApp.Group("/api", logger.New(), getJwtFunction(config))

	{
		api = api.Group("/current", apiKeyAuth(config))
		api.Get("/hello/:name", sayHello)
		companyGrouping(api, config)
		restaurantGrouping(api, config)
	}

}

func simulatedLogin(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.FormValue("user")
		pass := c.FormValue("pass")

		// Throws Unauthorized error
		if user != config.SimulatedUser || pass != config.SimulatedPassword {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name":  "John Doe",
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.OAuthSecret))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}

func companyGrouping(api fiber.Router, config *app.Configuration) {
	api.Get("/company/:companyId", func(c *fiber.Ctx) error {

		if id, err := strconv.Atoi(c.Params("companyId", "0")); err == nil && id > 0 {
			var comp models.Company
			repo := helpers.MainRepository{DB: config.DB}
			if _, err := repo.GetCompanyByID(&comp, uint(id)); err == nil {
				return c.Status(fiber.StatusOK).JSON(&comp)
			}

		}
		return c.SendStatus(fiber.StatusNotFound)
	})

}

func restaurantGrouping(group fiber.Router, config *app.Configuration) {
	api := group.Group("/restaurant")

	// TODO need to better understand gofiber and build out a better way to DRY the id validation
	//api := group.Group("/restaurant", func(c *fiber.Ctx) error {
	//	if c.Params("restaurantID", "0") != "0"{
	//		return c.Next()
	//	}
	//	return c.SendStatus(fiber.StatusNotFound)
	//})

	api.Get("/:restaurantID/menu", func(c *fiber.Ctx) error {
		if id, err := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64); err == nil && id > 0 {
			var menus []models.Menu
			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetAllMenusByRestaurantId(&menus, id); err == nil {
				return c.Status(fiber.StatusOK).JSON(&menus)
			}
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Post("/:restaurantID/menu", func(c *fiber.Ctx) error {
		if id, err := strconv.Atoi(c.Params("restaurantID", "0")); err == nil && id > 0 {

			// Can we even parse this object?
			var menu models.Menu
			if parseErr := c.BodyParser(&menu); parseErr != nil {
				return c.Status(fiber.StatusBadRequest).SendString(parseErr.Error())
				//return c.SendStatus(fiber.StatusBadRequest)
			}
			// If we can find a restaurant that exists, we can add it
			repo := helpers.MainRepository{DB: config.DB}
			rId := uint(id)
			if exists := repo.CheckRestaurantIdExists(rId); exists {
				menu.RestaurantID = rId

				if err := repo.SaveMenu(&menu); err == nil {
					return c.Status(fiber.StatusCreated).JSON(&menu)
				}
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})

	api.Get("/:restaurantID/menu/:menuID", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		// Make sure we're getting valid data
		if rErr != nil || mErr != nil || rID <= 0 || mID <= 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		repo := helpers.MainRepository{DB: config.DB}

		var menu models.Menu
		err := repo.GetMenuByMenuAndRestaurantID(&menu, mID, rID)

		if menu.RestaurantID == uint(rID) && err == nil {
			return c.JSON(&menu)
		}

		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Delete("/:restaurantID/menu/:menuID", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		// Make sure we're getting valid data
		if rErr == nil || mErr == nil || rID >= 1 || mID >= 1 {
			repo := helpers.MainRepository{DB: config.DB}
			err := repo.DeleteMenuByRestaurantAndMenuId(mID, rID)
			if err == nil {
				return c.SendStatus(fiber.StatusNoContent)
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})

	api.Get("/:restaurantID/menu/:menuID/menuItem/", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && rID >= 1 && mID >= 1 {
			var items []models.MenuItem

			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetAllMenuItemsByRestaurantAndMenuAndMenuItemIDs(&items, mID, rID); err == nil {
				return c.JSON(items)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Get("/:restaurantID/menu/:menuID/menuItem/:itemID", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && iErr == nil &&
			rID >= 1 && mID >= 1 && iID >= 1 {
			var item models.MenuItem

			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetMenuItemByRestaurantAndMenuAndMenuItemIDs(&item, iID, mID, rID); err == nil {
				return c.JSON(item)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Put("/:restaurantID/menu/:menuID/menuItem/:itemID", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && iErr == nil && rID >= 1 && mID >= 1 && iID >= 1 {

			var item map[string]interface{}
			err := c.BodyParser(&item)

			if err == nil {
				if num, pErr := strconv.ParseUint(fmt.Sprintf("%#v", item["ID"]), 10, 64); pErr == nil && num == mID {

					repo := helpers.MainRepository{DB: config.DB}
					if err := repo.UpdateMenuItem(&item, num, mID, rID); err != nil {
						return c.SendStatus(fiber.StatusInternalServerError)
					}
					return c.SendStatus(fiber.StatusAccepted)
				}
			}
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Delete("/:restaurantID/menu/:menuID/menuItem/:itemID", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && iErr == nil && rID >= 1 && mID >= 1 && iID >= 1 {
			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.DeleteMenuItem(iID, mID, rID); err == nil {

				c.SendStatus(fiber.StatusAccepted)
			}

		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Post("/:restaurantID/menu/:menuID/menuItem", func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		// Do we have valid parameters?
		if rErr == nil && mErr == nil && rID >= 1 && mID >= 1 {

			var item models.MenuItem
			// Can we even parse this object?
			if parseError := c.BodyParser(&item); parseError == nil {
				// We have to make sure that a restaurant and a corresponding menu exists
				repo := helpers.MainRepository{DB: config.DB}
				exists := repo.CheckMenuAndRestaurantIDsExist(mID, rID)
				if exists {
					item.MenuID = uint(mID)
					if count, err := repo.SaveMenuItem(&item); count <= 0 || err != nil {
						return c.SendStatus(fiber.StatusInternalServerError)
					}

					return c.SendStatus(fiber.StatusCreated)
				}
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	})

}

func menuItemGrouping(g fiber.Router, config *app.Configuration) {
	g.Get("MenuItems/:restaurantId", func(c *fiber.Ctx) error {
		if id, err := strconv.Atoi(c.Params("restaurantId", "0")); err == nil && id > 0 {

			//config.DB.Find
		}
		return c.Status(fiber.StatusNotFound).SendString("Cannot find requested page")
	})
}
