package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	UserRole  []UserRole `json:"userRole"`
	Addresses []Address  `json:"addresses"`
	//Roles     []UserRole `json:"roles" gorm:"foreignKey:ID""`
	//Addresses []Address  `json:"addresses"`
}
