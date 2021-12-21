package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string      `json:"name"`
	Price       int32       `json:"price"`
	Description string      `json:"description"`
	Images      []MenuImage `json:"images"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MenuID      uint        `json:"menuID"`
	//Restaurant  []*Restaurant `gorm:"foreignKey:ID"`
}

type MenuImage struct {
	gorm.Model
	ImageURL   string `json:"ImageURL"`
	MenuItemId uint   `json:"menuItemId"`
}

type Menu struct {
	gorm.Model
	Name         string     `json:"name"`
	Active       bool       `json:"active"`
	MenuItems    []MenuItem `json:"menuItems"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RestaurantID uint       `json:"restaurantID"`
}
