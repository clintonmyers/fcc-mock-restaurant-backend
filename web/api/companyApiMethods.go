package api

import (
	"github.com/clintonmyers/fcc-mock-restaurant-backend/app"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func getCompanyById(config *app.Configuration) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if id, err := strconv.Atoi(c.Params("companyId", "0")); err == nil && id > 0 {
			var comp models.Company
			//Repo := helpers.MainRepository{DB: config.DB}
			if _, err := Repo.GetCompanyByID(&comp, uint(id)); err == nil {
				return c.Status(fiber.StatusOK).JSON(&comp)
			}
		}

		return c.SendStatus(fiber.StatusNotFound)
	}

}
