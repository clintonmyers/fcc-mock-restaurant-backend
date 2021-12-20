package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string   `json:"name"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	//Restaurant  []*Restaurant `gorm:"foreignKey:ID"`
}

type Menu struct {
	gorm.Model
	RestaurantId uint   `json:"restaurantId"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
}
