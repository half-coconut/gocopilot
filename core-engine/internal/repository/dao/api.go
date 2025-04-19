package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type APIDAO interface {
	Insert(ctx context.Context, api API) (int64, error)
	UpdateById(ctx context.Context, api API) error
	FindByUId(ctx context.Context, id int64) ([]API, error)
	FindByAId(ctx context.Context, aid int64) (API, error)
}

type GORMAPIDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func (dao *GORMAPIDAO) FindByAId(ctx context.Context, aid int64) (API, error) {
	var api API
	err := dao.db.WithContext(ctx).Where("id=?", aid).Find(&api).Error
	return api, err
}

func (dao *GORMAPIDAO) FindByUId(ctx context.Context, id int64) ([]API, error) {
	var api []API
	err := dao.db.WithContext(ctx).Where("creator_id=?", id).Find(&api).Error
	return api, err
}

func NewAPIDAO(l logger.LoggerV1, db *gorm.DB) APIDAO {
	return &GORMAPIDAO{
		db: db,
		l:  l,
	}
}

func (dao *GORMAPIDAO) Insert(ctx context.Context, api API) (int64, error) {
	now := time.Now().UnixMilli()
	api.Ctime = now
	api.Utime = now
	api.UpdaterId = api.CreatorId
	err := dao.db.WithContext(ctx).Create(&api).Error
	return api.Id, err
}

func (dao *GORMAPIDAO) UpdateById(ctx context.Context, api API) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&api).Where("id=?", api.Id).
		Updates(map[string]interface{}{
			"name":         api.Name,
			"url":          api.URL,
			"params":       api.Params,
			"type":         api.Type,
			"body":         api.Body,
			"header":       api.Header,
			"method":       api.Method,
			"project":      api.Project,
			"debug_result": api.DebugResult,
			"utime":        now,
			"updater_id":   api.UpdaterId,
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

type API struct {
	Id          int64 `gorm:"primaryKey,autoIncrement"`
	Name        sql.NullString
	URL         sql.NullString
	Params      sql.NullString
	Type        sql.NullString
	Body        sql.NullString
	Header      sql.NullString
	Method      sql.NullString
	Project     sql.NullString
	DebugResult sql.NullString
	CreatorId   int64
	UpdaterId   int64
	Ctime       int64
	Utime       int64
}
