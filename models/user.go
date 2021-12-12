package models

import "gorm.io/gorm"

type User struct {
	//UID       uint `gorm:"primarykey"`
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt gorm.DeletedAt `gorm:"index"`
	gorm.Model
	Role      Role   `json:"role"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address
	RestaurantId uint
	//Restaurant *Restaurant
	//Address   string `json:"address"`
	//City      string `json:"city"`
	//State     string `json:"state"`
	//Zip       string `json:"zip"`
}

const RoleCustomer Role = 1
const RoleAdmin Role = 1

type Role int8

const (
	Customer Role = iota
	Admin
)

func (r Role) String() string {
	return []string{"customer", "admin"}[r]
}

//Users
//
//    id (primary key)
//
//    role (customer or restaurant admin)
//
//    username
//
//    firstName
//
//    lastName
//
//    address
//
//    city
//
//    state
//
//    zip
