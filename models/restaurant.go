package models

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model
	Address
	Users     []*User `gorm:"foreignKey:RestaurantId"`
	Company   Company
	CompanyId uint
}
