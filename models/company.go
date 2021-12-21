package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Restaurants []Restaurant `json:"restaurants"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
