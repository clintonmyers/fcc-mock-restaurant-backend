package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	//Restaurants []Restaurant `json:"restaurants"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	//Restaurants []Restaurant `json:"restaurants"gorm:"foreignKey:CompanyID"`
	Restaurants []Restaurant `json:"restaurant"`
}
