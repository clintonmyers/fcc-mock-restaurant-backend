package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	Role         string `json:"role"`
	RestaurantID uint   `json:"restaurantID"`
	UserId       uint   `json:"userId"`
	//Restaurant Restaurant `json:"restaurant"`
}
