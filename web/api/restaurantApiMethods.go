package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func getMenuByRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if id, err := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64); err == nil && id > 0 {
			var menus []models.Menu
			//Repo := helpers.MainRepository{DB: config.DB}

			if err := Repo.GetAllMenusByRestaurantId(&menus, id); err == nil {
				return c.Status(fiber.StatusOK).JSON(&menus)
			}
		}
		return c.SendStatus(fiber.StatusNotFound)
	}
}

func postMenuByRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if id, err := strconv.Atoi(c.Params("restaurantID", "0")); err == nil && id > 0 {

			// Can we even parse this object?
			var menu models.Menu
			if parseErr := c.BodyParser(&menu); parseErr != nil {
				return c.Status(fiber.StatusBadRequest).SendString(parseErr.Error())
				//return c.SendStatus(fiber.StatusBadRequest)
			}
			// If we can find a restaurant that exists, we can add it
			//Repo := helpers.MainRepository{DB: config.DB}
			rId := uint(id)
			if exists := Repo.CheckRestaurantIdExists(rId); exists {
				menu.RestaurantID = rId

				if err := Repo.SaveMenu(&menu); err == nil {
					return c.Status(fiber.StatusCreated).JSON(&menu)
				}
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func getMenuByIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		// Make sure we're getting valid data
		if rErr != nil || mErr != nil || rID <= 0 || mID <= 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		//Repo := helpers.MainRepository{DB: config.DB}

		var menu models.Menu
		err := Repo.GetMenuByMenuAndRestaurantID(&menu, mID, rID)

		if menu.RestaurantID == uint(rID) && err == nil {
			return c.JSON(&menu)
		}

		return c.SendStatus(fiber.StatusNotFound)
	}
}

func deleteMenuByIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		// Make sure we're getting valid data
		if rErr == nil || mErr == nil || rID >= 1 || mID >= 1 {
			//Repo := helpers.MainRepository{DB: config.DB}
			err := Repo.DeleteMenuByRestaurantAndMenuId(mID, rID)
			if err == nil {
				return c.SendStatus(fiber.StatusNoContent)
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}

func getMenuItemByMenuIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && rID >= 1 && mID >= 1 {
			var items []models.MenuItem

			//Repo := helpers.MainRepository{DB: config.DB}

			if err := Repo.GetAllMenuItemsByRestaurantAndMenuAndMenuItemIDs(&items, mID, rID); err == nil {
				return c.JSON(items)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)
	}
}
func getItemByIdAndMenuIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && iErr == nil &&
			rID >= 1 && mID >= 1 && iID >= 1 {
			var item models.MenuItem

			//Repo := helpers.MainRepository{DB: config.DB}

			if err := Repo.GetMenuItemByRestaurantAndMenuAndMenuItemIDs(&item, iID, mID, rID); err == nil {
				return c.JSON(item)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)
	}
}

func putItemByIdAndMenuIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if (rErr == nil && mErr == nil && iErr == nil) && (rID >= 1 && mID >= 1 && iID >= 1) {

			var item map[string]interface{}
			err := c.BodyParser(&item)

			if err == nil {
				if num, pErr := strconv.ParseUint(fmt.Sprintf("%#v", item["ID"]), 10, 64); pErr == nil && num == mID {

					//Repo := helpers.MainRepository{DB: config.DB}
					if err := Repo.UpdateMenuItem(&item, num, mID, rID); err != nil {
						return c.SendStatus(fiber.StatusInternalServerError)
					}
					return c.SendStatus(fiber.StatusAccepted)
				}
			}
		}
		return c.SendStatus(fiber.StatusNotFound)
	}
}

func deleteItemByIdAndMenuIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)
		iID, iErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		if rErr == nil && mErr == nil && iErr == nil && rID >= 1 && mID >= 1 && iID >= 1 {
			//Repo := helpers.MainRepository{DB: config.DB}

			if err := Repo.DeleteMenuItem(iID, mID, rID); err == nil {

				c.SendStatus(fiber.StatusAccepted)
			}

		}
		return c.SendStatus(fiber.StatusNotFound)
	}
}
func postMenuItemByItemIdAndRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rID, rErr := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64)
		mID, mErr := strconv.ParseUint(c.Params("menuID", "0"), 10, 64)

		// Do we have valid parameters?
		if rErr == nil && mErr == nil && rID >= 1 && mID >= 1 {

			var item models.MenuItem
			// Can we even parse this object?
			if parseError := c.BodyParser(&item); parseError == nil {
				// We have to make sure that a restaurant and a corresponding menu exists
				//Repo := helpers.MainRepository{DB: config.DB}
				exists := Repo.CheckMenuAndRestaurantIDsExist(mID, rID)
				if exists {
					item.MenuID = uint(mID)
					if count, err := Repo.SaveMenuItem(&item); count <= 0 || err != nil {
						return c.SendStatus(fiber.StatusInternalServerError)
					}

					return c.SendStatus(fiber.StatusCreated)
				}
			}
		}
		return c.SendStatus(fiber.StatusBadRequest)
	}
}
