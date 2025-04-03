package dao

import (
	"TestCopilot/TestEngine/internal/repository/dao/note"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&note.Note{},
		&note.PublishedNote{},
		&API{},
		&Task{},

		&Interactive{},
		&UserLikeBiz{},
		&Collection{},
		&UserCollectionBiz{},
		&Job{},
	)
}
