package app

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Configuration struct {
	DB                  *gorm.DB
	WebApp              *fiber.App
	MaxIdle             int
	MaxOpenConn         int
	LifetimeMinutes     int
	Production          bool
	Port                string
	LocalDB             string
	DatabaseUrl         string
	GenerateTestData    bool
	ApiKey              string
	AutoMigrate         bool
	DeleteLocalDatabase bool
	GenerateLocalData   bool
}
