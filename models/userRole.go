package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	Role         string `json:"role"`
	RestaurantId uint   `json:"restaurantId"`
	UserId       uint   `json:"userId"`
	//Restaurant Restaurant `json:"restaurant"`
}
