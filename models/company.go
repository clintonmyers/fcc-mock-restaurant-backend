package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Company struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	//Restaurants []Restaurant `json:"restaurants"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	//Restaurants []Restaurant `json:"restaurants"gorm:"foreignKey:CompanyID"`
	Restaurants []Restaurant `json:"restaurant"gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CompanyID"`
}

func (b *Company) BeforeCreate(tx *gorm.DB) (err error) {
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
