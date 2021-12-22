package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

func (m *MainRepository) GetCompanyByID(c *models.Company, id uint) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid companyID")
	}
	m.DB.Preload(clause.Associations).Preload("Restaurants.Testimonials."+clause.Associations).Preload("Restaurants.Menus.MenuItems."+clause.Associations).Find(c, id)
	return 0, nil
}
