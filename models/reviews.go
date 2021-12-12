package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Restaurant *Restaurant `gorm:"foreignKey:ID"`
	Header     string
	Body       string
	Images     []string `gorm:"type:text[]"`
	//https://stackoverflow.com/questions/64035165/unsupported-data-type-error-on-gorm-field-where-custom-valuer-returns-nil
}
