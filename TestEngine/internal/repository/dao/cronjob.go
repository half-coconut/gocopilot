package dao

import (
	"TestCopilot/TestEngine/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
)

type CronJobDAO interface {
	Preempt(ctx context.Context) (CronJob, error)
	Insert(ctx context.Context, job CronJob) (int64, error)
	UpdateById(ctx context.Context, job CronJob) error
	GetJobById(ctx context.Context, jid int64) (CronJob, error)
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
}

type GORMCronJobDAO struct {
	db *gorm.DB
	l  logger.LoggerV1
}

func NewGORMCronJobDAO(db *gorm.DB, l logger.LoggerV1) CronJobDAO {
	return &GORMCronJobDAO{
		db: db,
		l:  l}
}

func (dao *GORMCronJobDAO) GetJobById(ctx context.Context, jid int64) (CronJob, error) {
	var job CronJob
	err := dao.db.WithContext(ctx).Where("id=?", jid).Find(&job).Error
	return job, err
}

func (dao *GORMCronJobDAO) Preempt(ctx context.Context) (CronJob, error) {
	db := dao.db.WithContext(ctx)
	for {

		now := time.Now()
		var j CronJob
		// 在所有即将满足条件的 job, 获取 first
		err := db.Model(&Job{}).Where("status = ? AND next_time <= ?", cronjobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return CronJob{}, err
		}

		// 乐观锁，CAS 操作，compare and swap
		res := db.Where("id=? AND version=?", j.Id, j.Version).Model(&Job{}).
			Updates(map[string]any{
				"status":  cronjobStatusRunning,
				"version": j.Version + 1, // 乐观锁，用于并发控制
				"utime":   now,
			})
		if res.Error != nil {
			return CronJob{}, err
		}
		if res.RowsAffected == 0 {
			//return Job{}, errors.New("没抢到")
			// 继续下一轮
			continue
		}
		return j, nil
	}
}

func (dao *GORMCronJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return dao.db.WithContext(ctx).Model(&CronJob{}).Where("id=?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (dao *GORMCronJobDAO) Insert(ctx context.Context, job CronJob) (int64, error) {
	now := time.Now().UnixMilli()
	job.Ctime = now
	job.Utime = now
	job.Status = cronjobStatusWaiting
	err := dao.db.WithContext(ctx).Create(&job).Error
	return job.Id, err
}

func (dao *GORMCronJobDAO) UpdateById(ctx context.Context, job CronJob) error {
	now := time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&job).Where("id=?", job.Id).
		Updates(map[string]interface{}{
			"name":        job.Name,
			"description": job.Description,
			"type":        job.Type,
			"cron":        job.Cron,
			"http_cfg":    job.HttpCfg,
			"task_id":     job.TaskId,
			"duration":    job.Duration,
			"retry":       job.Retry,
			"max_retries": job.MaxRetries,
			"status":      cronjobStatusWaiting,
			"utime":       now,
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

type CronJob struct {
	Id          int64          `gorm:"primaryKey,autoIncrement"`
	Name        sql.NullString `gorm:"unique"`
	Description sql.NullString
	Type        sql.NullString
	Cron        sql.NullString
	HttpCfg     sql.NullString
	TaskId      int64
	Duration    int64
	Retry       bool
	MaxRetries  uint64

	Version  int
	Status   int
	NextTime int64 `gorm:"index"` // 更好的 next_time 和 status 的联合索引

	Executor string // 内部调用

	CreatorId int64
	Ctime     int64
	Utime     int64
}

const (
	cronjobStatusWaiting = 0
	cronjobStatusRunning = 1
	// 暂停调度
	cronjobStatusPaused = 2
)
