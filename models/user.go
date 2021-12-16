package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Role      Role   `json:"role"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address
	RestaurantId uint
}

const RoleCustomer Role = 0
const RoleAdmin Role = 1

type Role int8

const (
	Customer Role = iota
	Admin
)

func (r Role) String() string {
	return []string{"customer", "admin"}[r]
}
