package helpers

import (
	"errors"
	"github.com/clintonmyers/fcc-mock-restaurant-backend/models"
	"gorm.io/gorm/clause"
)

type TestimonialRepository interface {
	SaveTestimonial(*models.Testimonial) (int64, error)
	GetTestimonialById(*models.Testimonial, uint) error
	GetAllTestimonials(*[]models.Testimonial) error
	GetAllTestimonialsByRestaurantId(*[]models.Testimonial, uint) error
	GetAllTestimonialsByRestaurant(*[]models.Testimonial, *models.Restaurant) error
}

func (m *MainRepository) SaveTestimonial(testimonial *models.Testimonial) (int64, error) {
	return m.DB.Save(testimonial).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetTestimonialById(testimonial *models.Testimonial, u uint) error {
	if u <= 0 {
		return errors.New("invalid id")
	}
	m.DB.Preload(clause.Associations).Find(testimonial, u)

	return m.DB.Error
}

//func (m *MainRepository) GetTestimonialByUsername(testimonial *models.Testimonial, s string) error {
//	if len(s) == 0 {
//		return errors.New("invalid username")
//	}
//	m.DB
//	//TODO implement me
//	panic("implement me")
//}

func (m *MainRepository) GetAllTestimonials(testimonials *[]models.Testimonial) error {
	m.DB.Preload(clause.Associations).Find(testimonials)
	return m.DB.Error
}

func (m *MainRepository) GetAllTestimonialsByRestaurantId(testimonials *[]models.Testimonial, u uint) error {
	m.DB.Preload(clause.Associations).Where("restaurantId = ?", u).Find(testimonials)
	return m.DB.Error
}

func (m *MainRepository) GetAllTestimonialsByRestaurant(testimonials *[]models.Testimonial, restaurant *models.Restaurant) error {
	return m.GetAllTestimonialsByRestaurantId(testimonials, restaurant.ID)
}

func (m *MainRepository) GetAllTestimonialsByCompanyTestimonials(t []models.Testimonial, c *models.Company) error {
	//var testimonials *[]models.Testimonial
	var err error
	for _, restaurant := range c.Restaurants {
		testimonials := make([]models.Testimonial, 0)
		err = m.GetAllTestimonialsByRestaurant(&testimonials, &restaurant)
		if err != nil {
			return err
		}
		//t = append(t, &testimonials)
		for _, test := range testimonials {
			t = append(t, test)
		}
	}
	return err
}
