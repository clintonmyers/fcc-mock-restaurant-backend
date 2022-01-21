package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string        `json:"username"`
	SubId     string        `json:"-"` // This is the ID from the OIDC and we don't want that sent back and forth
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	UserRole  []UserRole    `json:"userRole"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Addresses []UserAddress `json:"addresses"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
