package helpers

import "github.com/clintonmyers/fcc-mock-restaurant-backend/models"

//https://stackoverflow.com/questions/66392372/select-exists-with-gorm
func (m *MainRepository) CheckRestaurantIdExists(u uint) bool {
	if u <= 0 {
		return false
	}
	//var exists bool
	//m.DB.Model(&models.Restaurant{}).Find(&exists, u)
	var exists bool
	db := m.DB
	err := db.Model(&models.Restaurant{}).
		Select("count(*) > 0").
		Where("id = ?", u).
		Find(&exists).
		Error
	return err == nil
}
