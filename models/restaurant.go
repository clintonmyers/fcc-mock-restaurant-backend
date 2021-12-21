package models

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	Addresses    []Address     `json:"addresses"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	Testimonials []Testimonial `json:"testimonials"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	Menus        []Menu        `json:"menus"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	CompanyID    uint          `json:"companyID"`
}
