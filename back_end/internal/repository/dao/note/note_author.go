package note

import (
	"context"
	"egg_yolk/pkg/logger"
	"errors"
	"gorm.io/gorm"
	"time"
)

type AuthorDAO interface {
	Insert(ctx context.Context, note Note) (int64, error)
	UpdateById(ctx context.Context, note Note) error
}
type GORMAuthorDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func NewNoteAuthorDAO(l logger.LoggerV1, db *gorm.DB) AuthorDAO {
	return &GORMAuthorDAO{
		db: db,
		l:  l,
	}
}

func (dao *GORMAuthorDAO) Insert(ctx context.Context, note Note) (int64, error) {
	now := time.Now().UnixMilli()
	note.Ctime = now
	note.Utime = now
	err := dao.db.WithContext(ctx).Create(&note).Error
	return note.Id, err
}
func (dao *GORMAuthorDAO) UpdateById(ctx context.Context, note Note) error {
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
