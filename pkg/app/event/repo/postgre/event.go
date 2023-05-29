package postgre

import (
	"fmt"
	"nostr-ex/pkg/db"
	"nostr-ex/pkg/models"

	"gorm.io/gorm"
)

func GetMainDB() (*gorm.DB, error) {
	db := db.GetMainDB()
	if db == nil {
		return nil, fmt.Errorf("main db is nil")
	}

	return db, nil
}

func SaveEvent(data *models.Event) error {
	db, err := GetMainDB()
	if err != nil {
		return err
	}

	dbc := db.Create(data)
	return dbc.Error
}

func GetEvent(limit int) []models.Event {
	var ret []models.Event
	db, err := GetMainDB()
	if err != nil {
		return ret
	}

	db.Find(ret).Limit(limit)
	return ret
}
