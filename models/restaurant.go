package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Restaurant struct {
	gorm.Model
	Addresses    []RestaurantAddress `json:"addresses"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	Testimonials []Testimonial       `json:"testimonials"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	Menus        []Menu              `json:"menus"gorm:"constraint:OnUpdate:CASCADE,OnSave:CASCADE,OnDelete:SET NULL;"`
	CompanyID    uint                `json:"companyID"`
}

func (b *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
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
