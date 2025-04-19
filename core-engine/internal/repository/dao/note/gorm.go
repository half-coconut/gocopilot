package note

import (
	"context"
	"errors"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type GORMNoteDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func (dao *GORMNoteDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Note, error) {
	var res []Note
	// 保持排序的稳定
	err := dao.db.WithContext(ctx).
		Where("utime<?", start.UnixMilli()).
		Order("utime DESC").
		Offset(offset).
		Limit(limit).
		Find(&res).Error
	return res, err
}

func (dao *GORMNoteDAO) GetById(ctx context.Context, id int64) (Note, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMNoteDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Note, error) {
	var arts []Note
	// SELECT * FROM XXX WHERE XX order by aaa
	// 在设计 order by 语句的时候，要注意让 order by 中的数据命中索引
	// SQL 优化的案例：早期的时候，
	// 我们的 order by 没有命中索引的，内存排序非常慢
	// 你的工作就是优化了这个查询，加进去了索引
	// author_id => author_id, utime 的联合索引
	err := dao.db.WithContext(ctx).Model(&Note{}).
		Where("author_id = ?", author).
		Offset(offset).
		Limit(limit).
		// 升序排序。 utime ASC
		// 混合排序
		// ctime ASC, utime desc
		Order("utime DESC").
		//Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		//	{Column: clause.Column{Name: "utime"}, Desc: true},
		//	{Column: clause.Column{Name: "ctime"}, Desc: false},
		//}}).
		Find(&arts).Error
	return arts, err
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

func (dao *GORMNoteDAO) GetPubById(ctx context.Context, id int64) (PublishedNote, error) {
	var pub PublishedNote
	err := dao.db.WithContext(ctx).
		Where("id = ?", id).
		First(&pub).Error
	return pub, err
}
