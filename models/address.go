package models

type Address struct {
	Address      string `json:"address"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Active       bool   `json:"active"`
	UserId       uint   `json:"userId"`
	RestaurantId uint   `json:"restaurantId"`
}
