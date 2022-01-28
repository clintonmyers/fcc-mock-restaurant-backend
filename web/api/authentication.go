package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shareed2k/goth_fiber"
	"time"
)

func setupOAuth(fiberApp *fiber.App, config *app.Configuration) {

	fiberApp.Get("/login/:provider", goth_fiber.BeginAuthHandler)
	fiberApp.Get("/auth/:provider/callback", loginCallback(config))

	if config.SimulateOAuth {
		fiberApp.Post("/login", simulatedLogin(config))

	} else {
		fiberApp.Get("/login", func(c *fiber.Ctx) error {
			c.Set("provider", "google")
			return c.Redirect("/login/google", fiber.StatusTemporaryRedirect)
		})
	}

}

// For both regular and simulated login we need to get an app user based upon the OIDC user we get
// So we'll need to lookup a user by their sub ID

func loginCallback(config *app.Configuration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(ctx)
		if err != nil {
			fmt.Println(err)
		}

		if store, err := config.Store.Get(ctx); err == nil {
			defer store.Save()
			store.Set("user", user)

		} else {
			fmt.Println(err)
		}
		var repoUser models.User
		Repo.GetUserBySubId(&repoUser, user.UserID)
		if repoUser.ID == 0 {
			// TODO
			// We still need to work to setup the different restaurants based upon the requesting site

			roles := make([]models.UserRole, 0, 1)
			roles = append(roles, models.UserRole{
				Role:         "guest",
				RestaurantID: 1,
			})

			models.UserFromGormUser(&repoUser, &user)
			repoUser.UserRole = []models.UserRole{
				models.UserRole{
					Role:         "user",
					RestaurantID: 1,
				},
			}
			count, err := Repo.SaveUser(&repoUser)
			if count != 1 {
				fmt.Println("Hmmm didn't save properly.. expected 1 got: ", count, "\t", user)
			}
			if err != nil {
				fmt.Println("ERROR: ", err, "\t", user)
			}
		}
		isAdmin := false
		for _, role := range repoUser.UserRole {
			if role.Role == "admin" {
				isAdmin = true
			}
		}
		claims := jwt.MapClaims{
			"name":  repoUser.Username,
			"admin": isAdmin,
			"roles": repoUser.UserRole,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if t, err := token.SignedString([]byte(config.OAuthSecret)); err == nil {
			return ctx.Redirect(config.AuthRedirect+"?jwt="+t, fiber.StatusTemporaryRedirect)
		} else {
			fmt.Println(err)
		}
		return ctx.SendStatus(fiber.StatusInternalServerError)
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
			"name":  config.SimulatedUser,
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		if t, err := token.SignedString([]byte(config.OAuthSecret)); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			return c.JSON(fiber.Map{"token": t})
		}
	}
}

func jwtAuth(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		m := c.Method()
		// Only allow idempotent methods without access token
		if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
			return c.Next()
		}
		user := c.Locals("user").(*jwt.Token)
		if !user.Valid {
			fmt.Println("Invalid user, returning 404")
			return c.Status(fiber.StatusNotFound).SendString("Cannot Find Requested Page")

		}
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		fmt.Println("found user: ", name)
		isAdmin := claims["admin"].(bool)
		if isAdmin {
			return c.Next()
		}
		//return c.SendString("Welcome " + name)

		return c.Status(fiber.StatusNotFound).SendString("Cannot Find Requested Page")
	}
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

func getJwtFunction(config *app.Configuration) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.OAuthSecret),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fmt.Println("jwt error handler called, returning 404 to user")
			return c.SendStatus(fiber.StatusNotFound)
		},
		Filter: httpMethodBasedFilter,
	})
}

func httpMethodBasedFilter(ctx *fiber.Ctx) bool {
	m := ctx.Method()
	if m == fiber.MethodGet || m == fiber.MethodHead || m == fiber.MethodConnect || m == fiber.MethodOptions {
		return true
	}
	return false
}
