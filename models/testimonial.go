package models

import (
	"gorm.io/gorm"
)

type Testimonial struct {
	gorm.Model
	Header       string             `json:"header"`
	Body         string             `json:"body"`
	ImageUrls    []TestimonialImage `json:"imageUrls" gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	RestaurantID uint               `json:"restaurantID"`
}

type TestimonialImage struct {
	gorm.Model
	URL           string `json:"URL"`
	TestimonialId uint   `json:"TestimonialId"`
}
