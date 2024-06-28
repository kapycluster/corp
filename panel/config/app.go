package config

import "gorm.io/gorm"

// App struct is passed around across the application
type App struct {
	// Config as defined in the config.yaml file
	Config *Config

	// DB is the database connection
	DB *gorm.DB
}
