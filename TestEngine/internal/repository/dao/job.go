package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id =?").Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).Updates(map[string]any{
		"next_time": next,
	}).Error
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id=?", id).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := g.db.WithContext(ctx)
	for {

		now := time.Now()
		var j Job
		err := db.Model(&Job{}).Where("status = ? AND next_time <= ?", jobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return Job{}, err
		}

		// 乐观锁，CAS 操作，compare and swap
		res := db.Where("id=? AND version=?", j.Id, j.Version).Model(&Job{}).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"version": j.Version + 1, // 乐观锁，用于并发控制
				"utime":   now,
			})
		if res.Error != nil {
			return Job{}, err
		}
		if res.RowsAffected == 0 {
			//return Job{}, errors.New("没抢到")
			// 继续下一轮
			continue
		}
		return j, nil
	}
}

type Job struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Cfg    string `json:"cfg"`
	Name   string `gorm:"unique"`
	Status int
	// 下一次被调度的时间
	// 更好的 next_time 和 status 的联合索引
	NextTime int64 `gorm:"index"`
	Cron     string
	Version  int
	Executor string

	Ctime int64
	Utime int64
}

const (
	jobStatusWaiting = 0
	jobStatusRunning = 1
	// 暂停调度
	jobStatusPaused = 2
)
