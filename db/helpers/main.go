package helpers

import (
	"gorm.io/gorm"
)

type MainRepository struct {
	DB *gorm.DB
}
