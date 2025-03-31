package dao

import (
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
)

type TaskDAO interface {
	Insert(ctx context.Context, task Task) (int64, error)
	UpdateById(ctx context.Context, task Task) error
	FindByUId(ctx context.Context, id int64) ([]Task, error)
	FindByTId(ctx context.Context, tid int64) (Task, error)
}

type GORMTaskDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func (dao GORMTaskDAO) Insert(ctx context.Context, task Task) (int64, error) {
	now := time.Now().UnixMilli()
	task.Ctime = now
	task.Utime = now
	task.UpdaterId = task.CreatorId
	err := dao.db.WithContext(ctx).Create(&task).Error
	return task.Id, err
}

func (dao GORMTaskDAO) UpdateById(ctx context.Context, task Task) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&task).Where("id=?", task.Id).
		Updates(map[string]interface{}{
			"name":        task.Name,
			"a_ids":       task.AIds,
			"durations":   task.Durations,
			"workers":     task.Workers,
			"max_workers": task.MaxWorkers,
			"rate":        task.Rate,
			"utime":       now,
			"updater_id":  task.UpdaterId,
		})
	// 注意这里的处理，通过 RowsAffected==0，得知更新失败
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("更新数据失败！")
	}
	return err
}

func (dao GORMTaskDAO) FindByUId(ctx context.Context, id int64) ([]Task, error) {
	var task []Task
	err := dao.db.WithContext(ctx).Where("creator_id=?", id).Find(&task).Error
	return task, err
}

func (dao GORMTaskDAO) FindByTId(ctx context.Context, tid int64) (Task, error) {
	var task Task
	err := dao.db.WithContext(ctx).Where("id=?", tid).Find(&task).Error
	return task, err
}

func NewGORMTaskDAO(db *gorm.DB, l logger.LoggerV1) TaskDAO {
	return &GORMTaskDAO{db: db, l: l}
}

type Task struct {
	Id   int64 `gorm:"primaryKey,autoIncrement"`
	Name sql.NullString
	AIds sql.NullString

	Durations  int64
	Workers    uint64
	MaxWorkers uint64
	Rate       float64

	CreatorId int64
	UpdaterId int64
	Ctime     int64
	Utime     int64
}
