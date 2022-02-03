package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

func (m *MainRepository) GetAllMenusByRestaurantId(menus *[]models.Menu, id uint64) error {
	if id <= 0 {
		return errors.New("invalid Restaurant ID")
	}
	return m.DB.
		Where("restaurant_id = ?", id).
		Preload("MenuItems." + clause.Associations).
		Find(menus).Error
}

func (m *MainRepository) SaveMenu(menu *models.Menu) error {
	if menu.RestaurantID == 0 {
		return errors.New("unable to save menu element without attached restaurantID")
	}

	return m.DB.Save(menu).Error
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

func (m *MainRepository) SaveMenuItem(item *models.MenuItem) (int64, error) {
	return m.DB.Save(item).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetMenuItemByRestaurantAndMenuAndMenuItemIDs(item *models.MenuItem, itemID, menuID, restID uint64) error {
	// If the restaurant and the Menu have a connection then we're looking at the right item
	if m.CheckMenuAndRestaurantIDsExist(menuID, restID) == false {
		return errors.New("unable to find an associated menu and restaurant")
	}

	return m.DB.
		Where("menu_id = ?", menuID).
		Find(item, itemID).Error
}

func (m *MainRepository) GetAllMenuItemsByRestaurantAndMenuAndMenuItemIDs(items *[]models.MenuItem, menuID, restID uint64) error {
	// If the restaurant and the Menu have a connection then we're looking at the right item
	if m.CheckMenuAndRestaurantIDsExist(menuID, restID) == false {
		return errors.New("unable to find an associated menu and restaurant")
	}

	return m.DB.
		Where("menu_id = ?", menuID).
		Find(items).Error
}

func (m *MainRepository) UpdateMenuItem(item *map[string]interface{}, ItemID, menuID, restID uint64) error {
	if m.CheckMenuAndRestaurantIDsExist(menuID, restID) == false {
		return errors.New("unable to find an associated menu and restaurant")
	}
	return m.DB.Model(&models.MenuItem{}).
		Where("menu_id = ?", menuID).
		Where("ID = ?", ItemID).
		Save(item).Error
}

func (m *MainRepository) DeleteMenuItem(ItemID, menuID, restID uint64) error {
	if m.CheckMenuAndRestaurantIDsExist(menuID, restID) == false {
		return errors.New("unable to find an associated menu and restaurant")
	}
	return m.DB.
		Where("menu_id = ?", menuID).
		Delete(&models.MenuItem{}, ItemID).Error
}
