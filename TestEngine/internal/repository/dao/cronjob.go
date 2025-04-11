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
	PreemptByJId(ctx context.Context, jid int64) (CronJob, error)
	Insert(ctx context.Context, job CronJob) (int64, error)
	UpdateById(ctx context.Context, job CronJob) error
	GetJobById(ctx context.Context, jid int64) (CronJob, error)
	GetJobStatusById(ctx context.Context, jid int64) (int, error)
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
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
func (dao *GORMCronJobDAO) GetJobStatusById(ctx context.Context, jid int64) (int, error) {
	var job CronJob
	err := dao.db.WithContext(ctx).Where("id=?", jid).Find(&job).Error
	return job.Status, err
}
func (dao *GORMCronJobDAO) GetJobById(ctx context.Context, jid int64) (CronJob, error) {
	var job CronJob
	err := dao.db.WithContext(ctx).Where("id=?", jid).Find(&job).Error
	return job, err
}

func (dao *GORMCronJobDAO) PreemptByJId(ctx context.Context, jid int64) (CronJob, error) {
	db := dao.db.WithContext(ctx)
	for {
		now := time.Now().UnixMilli()
		var j CronJob
		err := db.Model(&CronJob{}).Where("id = ? AND status = ? AND next_time <= ?", jid, cronjobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return CronJob{}, err
		}
		// 使用乐观锁
		res := db.Where("id =? AND version=?", jid, j.Version).Model(&CronJob{}).
			Updates(map[string]any{
				"status":  cronjobStatusRunning,
				"version": j.Version + 1,
				"utime":   now,
			})
		if res.Error != nil {
			return CronJob{}, err
		}
		if res.RowsAffected == 0 {
			return CronJob{}, errors.New("没抢到")
		}
		return j, err
	}

}

func (dao *GORMCronJobDAO) Preempt(ctx context.Context) (CronJob, error) {
	db := dao.db.WithContext(ctx)
	// for 循环通过允许重试来提高了任务抢占的成功率
	for {

		now := time.Now()
		var j CronJob
		// 在所有即将满足条件的 job, 获取 first
		err := db.Model(&CronJob{}).Where("status = ? AND next_time <= ?", cronjobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return CronJob{}, err
		}

		// 乐观锁，CAS 操作，compare and swap
		res := db.Where("id=? AND version=?", j.Id, j.Version).Model(&CronJob{}).
			Updates(map[string]any{
				"status":  cronjobStatusRunning,
				"version": j.Version + 1, // 乐观锁，用于并发控制
				"utime":   now.UnixMilli(),
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

func (dao *GORMCronJobDAO) Release(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&CronJob{}).Where("id=?", id).Updates(map[string]any{
		"utime":  time.Now().UnixMilli(),
		"status": cronjobStatusWaiting,
	}).Error
}

func (dao *GORMCronJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return dao.db.WithContext(ctx).Model(&CronJob{}).Where("id=?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}
func (dao *GORMCronJobDAO) Stop(ctx context.Context, id int64) error {
	return dao.db.WithContext(ctx).Model(&CronJob{}).Where("id=?", id).Updates(map[string]any{
		"status": cronjobStatusPaused,
		"utime":  time.Now().UnixMilli(),
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
	// 等待，准备进入
	cronjobStatusWaiting = iota
	// 执行中
	cronjobStatusRunning
	// 暂停调度
	cronjobStatusPaused
)
