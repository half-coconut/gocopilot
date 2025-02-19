package dao

import (
	"egg_yolk/internal/repository/dao/note"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &note.Note{}, &note.PublishedNote{}, &API{})
}
