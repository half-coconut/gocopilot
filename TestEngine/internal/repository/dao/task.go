package dao

import (
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type TaskDAO interface {
	Insert(ctx context.Context, task Task) (int64, error)
	UpdateById(ctx context.Context, task Task) error
	FindByUId(ctx context.Context, id int64) ([]API, error)
	FindByTId(ctx context.Context, tid int64) (API, error)
}

type GORMTASKDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func (G GORMTASKDAO) Insert(ctx context.Context, task Task) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (G GORMTASKDAO) UpdateById(ctx context.Context, task Task) error {
	//TODO implement me
	panic("implement me")
}

func (G GORMTASKDAO) FindByUId(ctx context.Context, id int64) ([]API, error) {
	//TODO implement me
	panic("implement me")
}

func (G GORMTASKDAO) FindByTId(ctx context.Context, tid int64) (API, error) {
	//TODO implement me
	panic("implement me")
}

func NewGORMTaskDAO(db *gorm.DB, l logger.LoggerV1) TaskDAO {
	return &GORMTASKDAO{db: db, l: l}
}

type Task struct {
	Id   int64 `gorm:"primaryKey,autoIncrement"`
	Name sql.NullString
	APIs sql.NullString

	Durations  int64
	Workers    uint64 `json:"workers"`    // 并发数
	MaxWorkers uint64 `json:"maxWorkers"` // 最大持续时间
	Timeout    int64

	CreatorId int64
	UpdaterId int64
	Ctime     int64
	Utime     int64
}
