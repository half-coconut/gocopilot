package note

import (
	"context"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, note PublishedNote) error
}

type GORMReaderDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func NewNoteReaderDAO(l logger.LoggerV1, db *gorm.DB) ReaderDAO {
	return &GORMReaderDAO{
		db: db,
		l:  l,
	}
}

func (dao GORMReaderDAO) Upsert(ctx context.Context, note PublishedNote) error {
	now := time.Now().UnixMilli()
	note.Ctime = now
	note.Utime = now
	err := dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // 指定唯一索引或主键作为冲突条件
		//DoUpdates: clause.AssignmentColumns([]string{"title", "content", "utime"}),
		DoUpdates: clause.Assignments(map[string]interface{}{ // 指定需要更新的字段
			"title":   note.Title,
			"content": note.Content,
			"status":  note.Status,
			"utime":   now,
		}),
	}).Create(&note).Error
	return err
}
