package service

import (
	"TestCopilot/TestEngine/internal/domain"
	svcmocks "TestCopilot/TestEngine/internal/service/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestBatchRankingService_TopN(t *testing.T) {
	now := time.Now()
	teatCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) (NoteService, InteractiveService)
		wantErr   error
		wantNotes []domain.Note
	}{
		{
			name: "计算成功",
			// 模拟数据
			mock: func(ctrl *gomock.Controller) (NoteService, InteractiveService) {
				noteSvc := svcmocks.NewMockNoteService(ctrl)
				// 一批就搞完
				noteSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 3).
					Return([]domain.Note{
						{Id: 1, Utime: now, Ctime: now},
						{Id: 2, Utime: now, Ctime: now},
						{Id: 3, Utime: now, Ctime: now},
					}, nil)
				noteSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 3, 3).
					Return([]domain.Note{}, nil)
				intrSvc := svcmocks.NewMockInteractiveService(ctrl)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					"note", []int64{1, 2, 3}).
					Return(map[int64]domain.Interactive{
						1: {BizId: 1, LikeCnt: 1},
						2: {BizId: 2, LikeCnt: 2},
						3: {BizId: 3, LikeCnt: 3},
					}, nil)
				intrSvc.EXPECT().GetByIds(gomock.Any(),
					"note", []int64{}).
					Return(map[int64]domain.Interactive{}, nil)
				return noteSvc, intrSvc
			},
			wantNotes: []domain.Note{
				{Id: 3, Utime: now, Ctime: now},
				{Id: 2, Utime: now, Ctime: now},
				{Id: 1, Utime: now, Ctime: now},
			},
		},
	}

	for _, tc := range teatCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			noteSvc, intrSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(noteSvc, intrSvc).(*BatchRankingService)
			// 用于测试
			svc.batchSize = 3
			svc.n = 3
			svc.scoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			note, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantNotes, note)
		})
	}
}
