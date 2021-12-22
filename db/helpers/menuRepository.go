package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

func (m *MainRepository) GetALlMenusByRestaurantId(menus *[]models.Menu, id uint64) (int64, error) {
	if id <= 0 {
		return 0, errors.New("invalid Restaurant ID")
	}
	return m.DB.
		Where("restaurant_id = ?", id).
		Preload("MenuItems." + clause.Associations).
		Find(menus).RowsAffected, m.DB.Error
}

func (m *MainRepository) SaveMenu(menu *models.Menu) (int64, error) {
	if menu.RestaurantID == 0 {
		return 0, errors.New("unable to save menu element without attached restaurantID")
	}

	return m.DB.Save(menu).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetMenuByMenuAndRestaurantID(menu *models.Menu, menuId, restId uint64) error {
	if menuId <= 0 {
		return errors.New("invalid menuID")
	}
	return m.DB.Where("restaurant_id = ?", restId).
		//Preload(clause.Associations). // We dont' need this one
		Preload("MenuItems."+clause.Associations).
		Find(menu, menuId).
		Error

}

func (m MainRepository) DeleteMenuByRestaurantAndMenuId(menuId, restId uint64) error {
	if menuId <= 0 || restId <= 0 {
		return errors.New("invalid id's used")
	}

	return m.DB.Where("restaurant_id = ?", restId).
		Delete(&models.Menu{}, menuId).
		Error
}

/*

func (m *MainRepository) GetAllTestimonialsByRestaurantId(testimonials *[]models.Testimonial, u uint64) error {
	m.DB.Preload(clause.Associations).Where("restaurantId = ?", u).Find(testimonials)
	return m.DB.Error
}
*/
// https: //stackoverflow.com/questions/66392372/select-exists-with-gorm
func (m *MainRepository) CheckMenuAndRestaurantIDsExist(menuId, restId uint64) bool {
	if menuId <= 0 {
		return false
	}
	var exists bool
	db := m.DB
	err := db.Model(&models.Menu{}).
		Select("count(*) > 0").
		Where("restaurant_id = ?", restId).
		Find(&exists, menuId).
		Error
	return err == nil
}

func (m *MainRepository) SaveMenuItem(item *models.MenuItem) {
	m.DB.Save(item)
}
