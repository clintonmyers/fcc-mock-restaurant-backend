package api

import (
	"fmt"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/db/helpers"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func getMenuByRestaurantId(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if id, err := strconv.ParseUint(c.Params("restaurantID", "0"), 10, 64); err == nil && id > 0 {
			var menus []models.Menu
			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetAllMenusByRestaurantId(&menus, id); err == nil {
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
		repo := helpers.MainRepository{DB: config.DB}

		var menu models.Menu
		err := repo.GetMenuByMenuAndRestaurantID(&menu, mID, rID)

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
			repo := helpers.MainRepository{DB: config.DB}
			err := repo.DeleteMenuByRestaurantAndMenuId(mID, rID)
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

			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetAllMenuItemsByRestaurantAndMenuAndMenuItemIDs(&items, mID, rID); err == nil {
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

			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.GetMenuItemByRestaurantAndMenuAndMenuItemIDs(&item, iID, mID, rID); err == nil {
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

					repo := helpers.MainRepository{DB: config.DB}
					if err := repo.UpdateMenuItem(&item, num, mID, rID); err != nil {
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
			repo := helpers.MainRepository{DB: config.DB}

			if err := repo.DeleteMenuItem(iID, mID, rID); err == nil {

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
	}
}
