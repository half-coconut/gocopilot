package service

import (
	"TestCopilot/TestEngine/internal/domain"
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"log"
	"math"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
}

type BatchRankingService struct {
	noteSvc   NoteService
	intrSvc   InteractiveService
	batchSize int
	// 优先队列的容量
	n int
	// scoreFunc 用于测试，不能返回负数
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func NewBatchRankingService(noteSvc NoteService, intrSvc InteractiveService) *BatchRankingService {
	return &BatchRankingService{
		noteSvc:   noteSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			return float64(likeCnt-1) / math.Pow(float64(likeCnt+2), 1.5)
		},
	}
}

// 准备分批
func (svc *BatchRankingService) TopN(ctx context.Context) error {
	notes, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	// 在这里，存起来
	log.Println(notes)
	return nil
}

func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Note, error) {
	now := time.Now()
	offset := 0
	type Score struct {
		note  domain.Note
		score float64
	}
	// 注意这里使用优先队列，按照热度排序的小顶堆，顶点是最小的
	topN := queue.NewPriorityQueue[Score](svc.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}

		})
	for {
		notes, err := svc.noteSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Note, int64](notes,
			func(idx int, src domain.Note) int64 {
				return src.Id
			})
		// 找对应的点赞数据
		intrs, err := svc.intrSvc.GetByIds(ctx, "note", ids)
		if err != nil {
			return nil, err
		}
		// 合并计算 score
		// 排序
		for _, note := range notes {
			intr := intrs[note.Id]
			//if !ok{
			//	// 没有数据
			//	continue
			//}
			score := svc.scoreFunc(note.Utime, intr.LikeCnt)
			// 拿到热度最低的
			err = topN.Enqueue(Score{
				note:  note,
				score: score,
			})
			if err == queue.ErrOutOfCapacity {
				// 要求 topN 已经满了
				val, _ := topN.Dequeue()
				if val.score < score {
					err = topN.Enqueue(Score{
						note:  note,
						score: score,
					})
				}
			}
		}
		// 一批已经处理完了
		if len(notes) < svc.batchSize {
			// 这一批没取够
			break
		}
		// 这边要更新 offset
		offset += len(notes)
	}
	// 最后得出结论
	res := make([]domain.Note, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		res[i] = val.note
	}

	return res, nil
}
