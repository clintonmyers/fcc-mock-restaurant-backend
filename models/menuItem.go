package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string
	Price       float32
	Description string
	Image       string
	Restaurant  []*Restaurant `gorm:"foreignKey:ID"`
}
