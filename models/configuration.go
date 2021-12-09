package models

import "gorm.io/gorm"

type Configuration struct {
	DB         *gorm.DB
	Production bool
}
