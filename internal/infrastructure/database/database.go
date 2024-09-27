package database

import (
	"os"

	"github.com/gomarchy/estimate/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to the database")
	}

	DB = db
}

func Migrate() {
	DB.AutoMigrate(&domain.Breakout{})
	DB.AutoMigrate(&domain.User{})
}
