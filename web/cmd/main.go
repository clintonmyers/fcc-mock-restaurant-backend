package main

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/web/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"log"
)

func main() {
	fmt.Println("Hello world!")

	app := fiber.New()

	// https://github.com/gofiber/fiber/tree/master/middleware/pprof
	// pprof logging is underneath /debug/pprof
	app.Use(pprof.New())

	// https://docs.gofiber.io/api/middleware/cors
	// The default already has '*' as the allowed origins
	app.Use(cors.New())

	// GET /john
	//app.Get("/:name", func(c *fiber.Ctx) error {
	//	msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
	//	return c.SendString(msg) // => Hello john 👋!
	//})
	api.GetRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
