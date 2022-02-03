package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RestaurantAddress struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Active  bool   `json:"active"`
	//UserId       uint   `json:"userId"`
	RestaurantId uint `json:"restaurantId"`
}

type UserAddress struct {
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Active  bool   `json:"active"`
	UserId  uint   `json:"userId"`
	//RestaurantId uint `json:"restaurantId"`
}

func (b *RestaurantAddress) BeforeCreate(tx *gorm.DB) (err error) {
	cols := []clause.Column{}
	colsNames := []string{}
	for _, field := range tx.Statement.Schema.PrimaryFields {
		cols = append(cols, clause.Column{Name: field.DBName})
		colsNames = append(colsNames, field.DBName)
	}
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.AssignmentColumns(colsNames),
		DoNothing: true,
	})
	return nil
}
