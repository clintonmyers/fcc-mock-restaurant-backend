package models

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	Addresses    []Address     `json:"addresses"`
	CompanyId    uint          `json:"companyId"`
	Testimonials []Testimonial `json:"testimonials"`
}
