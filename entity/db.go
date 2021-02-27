package entity

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB initializes the database
func InitDB() error {
	if db != nil {
		return nil
	}

	d, err := gorm.Open(sqlite.Open("./dbdata/data.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = d

	db.AutoMigrate(&User{})

	return nil
}

// Cleanup cleans up any db resources
func Cleanup() {
	// does nothing right now
}
