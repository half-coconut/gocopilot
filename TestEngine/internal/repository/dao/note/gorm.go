package note

import (
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GORMNoteDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func (dao *GORMNoteDAO) Sync(ctx context.Context, note Note) (int64, error) {
	//TODO implement me
	panic("implement me")
}

// SyncStatus 需要开启事务，将线上库和制作库，均改为自己可见
func (dao *GORMNoteDAO) SyncStatus(ctx context.Context, id, authorId int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Note{}).Where("id = ? AND author_id =?", id, authorId).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 介入监控，例如prometheus只要频繁出现，就告警，然后手工介入排查
			return fmt.Errorf("id 是错的，或者作者不对，非自己文章，uid: %d，note_id: %d", authorId, id)
		}
		return tx.Model(&PublishedNote{}).Where("id = ?", id).
			Updates(map[string]interface{}{
				"status": status,
				"utime":  now,
			}).Error
	})
}

func NewNoteDAO(l logger.LoggerV1, db *gorm.DB) NoteDAO {
	return &GORMNoteDAO{
		db: db,
		l:  l,
	}
}

func (dao *GORMNoteDAO) Insert(ctx context.Context, note Note) (int64, error) {
	now := time.Now().UnixMilli()
	note.Ctime = now
	note.Utime = now
	err := dao.db.WithContext(ctx).Create(&note).Error
	return note.Id, err
}
func (dao *GORMNoteDAO) UpdateById(ctx context.Context, note Note) error {
	now := time.Now().UnixMilli()
	// 不允许修改 author_id，所以如果找不到，就返回报错
	res := dao.db.WithContext(ctx).Model(&note).Where("id=?", note.Id).
		Where("author_id", note.AuthorId).
		Updates(map[string]interface{}{
			"title":   note.Title,
			"content": note.Content,
			"status":  note.Status,
			"utime":   now,
		})
	// 注意这里的处理，通过 RowsAffected==0，得知更新失败
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败！")
	}
	return nil
}
