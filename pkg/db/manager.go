package db

import (
	"gorm.io/gorm"
)

var mainDB *gorm.DB

func SetMainDB(db *gorm.DB) {
	mainDB = db
}

func GetMainDB() *gorm.DB {
	return mainDB
}
