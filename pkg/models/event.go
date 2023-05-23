package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	SubID string `gorm:"type:char(64)"`
	Data  string `gorm:"type:varchar(512)"`
}
