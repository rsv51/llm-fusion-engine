package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and migrates the schemas.
func InitDatabase(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&User{}, &ProxyKey{}, &Group{}, &Provider{}, &ApiKey{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}