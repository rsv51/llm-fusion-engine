package database

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	sqlite "github.com/glebarez/sqlite"
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
	err = DB.AutoMigrate(&User{}, &ProxyKey{}, &Group{}, &Provider{}, &ApiKey{}, &RequestLog{}, &Model{}, &ModelMapping{})
	if err != nil {
		return nil, err
	}

	// Create default admin user if not exists
	if err := createDefaultAdminUser(DB); err != nil {
		return nil, err
	}

	return DB, nil
}

// createDefaultAdminUser creates a default admin user if no users exist.
func createDefaultAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return err
	}

	// Only create default user if no users exist
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		defaultUser := User{
			Username: "admin",
			Password: string(hashedPassword),
			IsAdmin:  true,
		}

		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}
	}

	return nil
}