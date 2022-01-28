package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/gofiber/fiber/v2"
)

func sayHello(c *fiber.Ctx) error {
	msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
	return c.SendString(msg)
}

func SetupRoutes(fiberApp *fiber.App, config *app.Configuration) {
	fmt.Println("Setting Routes")

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		store, err := config.Store.Get(c)
		if err != nil {
			return c.SendString(err.Error())
		}
		user := store.Get("user")
		fmt.Println(c.Cookies("testName", "??"))
		return c.JSON(user)

	})
	fiberApp.Get("/hello/:name", func(c *fiber.Ctx) error {
		fmt.Println("HELLO")
		sess, err := config.Store.Get(c)
		if err != nil {
			fmt.Println(err)
			return c.SendString(err.Error())
		}
		defer sess.Save()
		sess.Set("test", "TESTING")

		return c.SendString("TESTING OUTPUT")

	})

	setupOAuth(fiberApp, config)

	api := fiberApp.Group("/api", getJwtFunction(config))

	{
		//api = api.Group("/current", apiKeyAuth(config))
		api = api.Group("/current", jwtAuth(config))
		api.Get("/hello/:name", sayHello)

		companyGrouping(api, config)
		restaurantGrouping(api, config)
	}

}

func companyGrouping(api fiber.Router, config *app.Configuration) {
	api.Get("/company/:companyId", getCompanyById(config))
}

func restaurantGrouping(group fiber.Router, config *app.Configuration) {
	//SetRestaurantApiRepo(helpers.MainRepository{DB: config.DB})

	api := group.Group("/restaurant")

	api.Get("/:restaurantID/menu", getMenuByRestaurantId(config))
	api.Post("/:restaurantID/menu", postMenuByRestaurantId(config))

	api.Get("/:restaurantID/menu/:menuID", getMenuByIdAndRestaurantId(config))
	api.Delete("/:restaurantID/menu/:menuID", deleteMenuByIdAndRestaurantId(config))

	api.Get("/:restaurantID/menu/:menuID/menuItem/", getMenuItemByMenuIdAndRestaurantId(config))
	api.Post("/:restaurantID/menu/:menuID/menuItem", postMenuItemByItemIdAndRestaurantId(config))

	api.Get("/:restaurantID/menu/:menuID/menuItem/:itemID", getItemByIdAndMenuIdAndRestaurantId(config))
	api.Put("/:restaurantID/menu/:menuID/menuItem/:itemID", putItemByIdAndMenuIdAndRestaurantId(config))
	api.Delete("/:restaurantID/menu/:menuID/menuItem/:itemID", deleteItemByIdAndMenuIdAndRestaurantId(config))

}

//
//func menuItemGrouping(g fiber.Router, config *app.Configuration) {
//	g.Get("MenuItems/:restaurantId", func(c *fiber.Ctx) error {
//		if id, err := strconv.Atoi(c.Params("restaurantId", "0")); err == nil && id > 0 {
//
//			//config.DB.Find
//		}
//		return c.Status(fiber.StatusNotFound).SendString("Cannot find requested page")
//	})
//}
