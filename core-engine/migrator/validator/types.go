package validator

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/half-coconut/gocopilot/core-engine/migrator/events"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/half-coconut/gocopilot/core-engine/pkg/migrator"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	l         logger.LoggerV1
	p         events.Producer
	direction string
	batchSize int
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, l logger.LoggerV1, p events.Producer, direction string) *Validator[T] {
	return &Validator[T]{base: base, target: target,
		l: l, p: p, direction: direction}
}

func (v *Validator[T]) Validate(ctx context.Context) {
	v.validateBaseToTarget(ctx)
	v.validateTargetToBase(ctx)
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) {
	offset := -1
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		// 一条条的查询，进来就更新 offset
		offset++
		var src T
		err := v.base.WithContext(dbCtx).Offset(offset).Order("id").First(&src).Error
		cancel()
		switch err {
		case nil:
			// 查到数据了
			var dst T
			err := v.target.Where("id=?", src.ID()).First(&dst).Error
			switch err {
			case nil:
				// 找到了，要开始比较了
				// reflect.DeepEqual(src, dst)
				if !src.CompareTo(dst) {
					v.notify(ctx, src.ID(), events.InconsistentEventTypeNEQ)
				}

			case gorm.ErrRecordNotFound:
				// target 里面少了数据
				v.notify(ctx, src.ID(), events.InconsistentEventTypeTargetMissing)
			default:
				v.l.Error("查询 target 失败")
				continue
			}

		case gorm.ErrRecordNotFound:
			// 说明比对结束了
			return
		default:
			// 数据库错误
			v.l.Error("校验数据查询 base 出错")
			continue
		}

	}
}

func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := -v.batchSize
	for {
		offset += v.batchSize
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.target.WithContext(dbCtx).Offset(offset).Limit(v.batchSize).
			Order("id").First(&dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			return
		}
		switch err {
		case gorm.ErrRecordNotFound:
			// 没数据了
			return
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err := v.base.Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				// 计算差集
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			default:
				continue
			}
		default:
			continue
		}
		if len(dstTs) < v.batchSize {
			return
		}
	}
}
func (v *Validator[T]) notifyBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.notify(ctx, id, events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(ctx context.Context, id int64, typ string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.p.ProducerInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Direction: v.direction,
		Type:      typ,
	})
	cancel()
	if err != nil {
		v.l.Error("发送消息失败", logger.Error(err))
	}
}
