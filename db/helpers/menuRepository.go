package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

func (m *MainRepository) GetALlMenusByRestaurantId(menus *[]models.Menu, u uint) (int64, error) {
	if u <= 0 {
		return 0, errors.New("invalid Restaurant ID")
	}
	return m.DB.Preload(clause.Associations).Where("restaurant_id = ?", u).Find(menus).RowsAffected, m.DB.Error
}

/*

func (m *MainRepository) GetAllTestimonialsByRestaurantId(testimonials *[]models.Testimonial, u uint) error {
	m.DB.Preload(clause.Associations).Where("restaurantId = ?", u).Find(testimonials)
	return m.DB.Error
}
*/
