package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name        string
	Description string
	Restaurants []*Restaurant `gorm:"foreignKey:CompanyId"`
}
