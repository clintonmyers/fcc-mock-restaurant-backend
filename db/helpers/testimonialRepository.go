package helpers

import "github.com/clintonmyers/fcc-mock-restaurant-backend/models"

type TestimonialRepository interface {
	SaveTestimonial(*models.Testimonial) (int64, error)
	GetTestimonialById(*models.Testimonial, uint) error
	GetTestimonialByUsername(*models.Testimonial, string) error
	GetAllTestimonials([]models.Testimonial) error
	GetAllTestimonialsByRestaurantId(*[]models.Testimonial, uint) error
	GetAllTestimonialsByRestaurant(*[]models.Testimonial, *models.Restaurant) error
}

func (m *MainRepository) SaveTestimonial(testimonial *models.Testimonial) (int64, error) {
	return m.DB.Save(testimonial).RowsAffected, m.DB.Error
}

func (m *MainRepository) GetTestimonialById(testimonial *models.Testimonial, u uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *MainRepository) GetTestimonialByUsername(testimonial *models.Testimonial, s string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MainRepository) GetAllTestimonials(testimonials []models.Testimonial) error {
	//TODO implement me
	panic("implement me")
}

func (m *MainRepository) GetAllTestimonialsByRestaurantId(testimonials *[]models.Testimonial, u uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *MainRepository) GetAllTestimonialsByRestaurant(testimonials *[]models.Testimonial, restaurant *models.Restaurant) error {
	//TODO implement me
	panic("implement me")
}
