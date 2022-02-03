package models

import "gorm.io/gorm"

type UserRole struct {
	gorm.Model
	Role         string `json:"role"`
	RestaurantID uint   `json:"restaurantID"`
	UserId       uint   `json:"userId"`
	//Restaurant Restaurant `json:"restaurant"`
}

func (ur *UserRole) TruncateRole() TruncatedUserRole {
	return TruncatedUserRole{
		ID:           ur.ID,
		Role:         ur.Role,
		RestaurantID: ur.RestaurantID,
		UserId:       ur.UserId,
	}
}

type TruncatedUserRole struct {
	ID           uint   `gorm:"primarykey"`
	Role         string `json:"role"`
	RestaurantID uint   `json:"restaurantID"`
	UserId       uint   `json:"userId"`
}
