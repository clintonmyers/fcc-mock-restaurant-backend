package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	UserRole  []UserRole `json:"userRole"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Addresses []Address  `json:"addresses"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
