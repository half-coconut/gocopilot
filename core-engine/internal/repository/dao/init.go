package dao

import (
	"github.com/half-coconut/gocopilot/core-engine/internal/repository/dao/note"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&note.Note{},
		&note.PublishedNote{},
		&API{},
		&Task{},
		&Job{},
		&CronJob{},
	)
}
