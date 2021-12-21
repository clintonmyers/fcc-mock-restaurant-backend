package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string      `json:"name"`
	Price       float32     `json:"price"`
	Description string      `json:"description"`
	Images      []MenuImage `json:"images"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//Restaurant  []*Restaurant `gorm:"foreignKey:ID"`
}

type MenuImage struct {
	gorm.Model
	URL        string `json:"URL"`
	MenuItemId uint   `json:"menuItemId"`
}

type Menu struct {
	gorm.Model
	RestaurantId uint   `json:"restaurantId"`
	Name         string `json:"name"`
	Active       bool   `json:"active"`
}
