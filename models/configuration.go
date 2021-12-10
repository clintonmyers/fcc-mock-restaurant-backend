package models

import "gorm.io/gorm"

type Configuration struct {
	DB              *gorm.DB
	MaxIdle         int
	MaxOpenConn     int
	LifetimeMinutes int
	Production      bool
	Port            string
	LocalDB         string
}
