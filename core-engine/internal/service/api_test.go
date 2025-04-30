package service

import (
	"context"
	"errors"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"github.com/half-coconut/gocopilot/core-engine/internal/repository"
	repomocks "github.com/half-coconut/gocopilot/core-engine/internal/repository/mocks"
	"github.com/half-coconut/gocopilot/core-engine/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestAPIService_Save(t *testing.T) {
	now := time.Now()
	testcases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) repository.APIRepository
		aid     int64
		api     domain.API
		wantId  int64
		wantErr error
	}{
		{name: "创建 API 成功",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), domain.API{
					Name:        "测试 API",
					URL:         "http://www.wanderful.com",
					Params:      "",
					Body:        "[]",
					Header:      "",
					Method:      "get",
					Type:        "http",
					Project:     "demo",
					Debug:       false,
					DebugResult: domain.DebugLog{},
					Creator: domain.Editor{
						Id: int64(66),
					},
					Updater: domain.Editor{
						Id: int64(66),
					},
					Ctime: now,
					Utime: now,
				}).Return(int64(123), nil)
				return repo
			},
			api: domain.API{
				Name:        "测试 API",
				URL:         "http://www.wanderful.com",
				Params:      "",
				Body:        "[]",
				Header:      "",
				Method:      "get",
				Type:        "http",
				Project:     "demo",
				Debug:       false,
				DebugResult: domain.DebugLog{},
				Creator: domain.Editor{
					Id: int64(123),
				},
				Updater: domain.Editor{
					Id: int64(123),
				},
				Ctime: now,
				Utime: now,
			},
			aid:     int64(123),
			wantId:  int64(123),
			wantErr: nil,
		},
		{name: "更新 API 成功",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().Update(gomock.Any(), domain.API{
					Id:          int64(123),
					Name:        "测试 API",
					URL:         "http://www.wanderful.com",
					Params:      "",
					Body:        "[]",
					Header:      "",
					Method:      "get",
					Type:        "http",
					Project:     "demo",
					Debug:       false,
					DebugResult: domain.DebugLog{},
					Creator: domain.Editor{
						Id: int64(123),
					},
					Updater: domain.Editor{
						Id: int64(66),
					},
					Ctime: now,
					Utime: now,
				}).Return(nil)
				return repo
			},
			api: domain.API{
				Id:          int64(123),
				Name:        "测试 API",
				URL:         "http://www.wanderful.com",
				Params:      "",
				Body:        "[]",
				Header:      "",
				Method:      "get",
				Type:        "http",
				Project:     "demo",
				Debug:       false,
				DebugResult: domain.DebugLog{},
				Creator: domain.Editor{
					Id: int64(123),
				},
				Updater: domain.Editor{
					Id: int64(123),
				},
				Ctime: now,
				Utime: now,
			},
			aid:     int64(123),
			wantId:  int64(123),
			wantErr: nil,
		},
		{name: "创建 API 失败",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().Create(gomock.Any(), domain.API{
					Name:        "测试 API",
					URL:         "http://www.wanderful.com",
					Params:      "",
					Body:        "[]",
					Header:      "",
					Method:      "get",
					Type:        "http",
					Project:     "demo",
					Debug:       false,
					DebugResult: domain.DebugLog{},
					Creator: domain.Editor{
						Id: int64(66),
					},
					Updater: domain.Editor{
						Id: int64(66),
					},
					Ctime: now,
					Utime: now,
				}).Return(int64(0), errors.New("创建失败"))
				return repo
			},
			api: domain.API{
				Name:        "测试 API",
				URL:         "http://www.wanderful.com",
				Params:      "",
				Body:        "[]",
				Header:      "",
				Method:      "get",
				Type:        "http",
				Project:     "demo",
				Debug:       false,
				DebugResult: domain.DebugLog{},
				Creator: domain.Editor{
					Id: int64(123),
				},
				Updater: domain.Editor{
					Id: int64(123),
				},
				Ctime: now,
				Utime: now,
			},
			aid:     int64(123),
			wantId:  int64(0),
			wantErr: errors.New("创建失败"),
		},
		{name: "更新 API 失败",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().Update(gomock.Any(), domain.API{
					Id:          int64(123),
					Name:        "测试 API",
					URL:         "http://www.wanderful.com",
					Params:      "",
					Body:        "[]",
					Header:      "",
					Method:      "get",
					Type:        "http",
					Project:     "demo",
					Debug:       false,
					DebugResult: domain.DebugLog{},
					Creator: domain.Editor{
						Id: int64(123),
					},
					Updater: domain.Editor{
						Id: int64(66),
					},
					Ctime: now,
					Utime: now,
				}).Return(errors.New("更新失败"))
				return repo
			},
			api: domain.API{
				Id:          int64(123),
				Name:        "测试 API",
				URL:         "http://www.wanderful.com",
				Params:      "",
				Body:        "[]",
				Header:      "",
				Method:      "get",
				Type:        "http",
				Project:     "demo",
				Debug:       false,
				DebugResult: domain.DebugLog{},
				Creator: domain.Editor{
					Id: int64(123),
				},
				Updater: domain.Editor{
					Id: int64(123),
				},
				Ctime: now,
				Utime: now,
			},
			aid:     int64(123),
			wantId:  int64(123),
			wantErr: errors.New("更新失败"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewAPIService(tc.mock(ctrl), &logger.NopLogger{})
			id, err := svc.Save(context.Background(), tc.api, int64(66))
			assert.Equal(t, tc.wantId, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestAPIService_Detail(t *testing.T) {

	testcases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) repository.APIRepository
		aid     int64
		api     domain.API
		wantAPI domain.API
		wantErr error
	}{
		{name: "获取 Detail 成功",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().FindByAId(gomock.Any(), int64(123)).Return(domain.API{
					Id:   int64(123),
					Name: "测试 API",
				}, nil)
				return repo
			},
			api: domain.API{
				Id:   int64(123),
				Name: "测试 API",
			},
			aid: int64(123),
			wantAPI: domain.API{
				Id:   int64(123),
				Name: "测试 API",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewAPIService(tc.mock(ctrl), &logger.NopLogger{})
			api, err := svc.Detail(context.Background(), tc.aid)
			assert.Equal(t, tc.wantAPI, api)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestAPIService_List(t *testing.T) {

	testcases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.APIRepository
		aid      int64
		apis     []domain.API
		wantAPIs []domain.API
		wantErr  error
	}{
		{name: "获取 List 成功",
			mock: func(ctrl *gomock.Controller) repository.APIRepository {
				repo := repomocks.NewMockAPIRepository(ctrl)
				repo.EXPECT().FindByUId(gomock.Any(), int64(123)).Return([]domain.API{
					{
						Id:   int64(123),
						Name: "测试 API",
					},
					{
						Id:   int64(125),
						Name: "测试 API 2",
					},
				}, nil)
				return repo
			},
			aid: int64(123),
			wantAPIs: []domain.API{
				{
					Id:   int64(123),
					Name: "测试 API",
				},
				{
					Id:   int64(125),
					Name: "测试 API 2",
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewAPIService(tc.mock(ctrl), &logger.NopLogger{})
			apis, err := svc.List(context.Background(), tc.aid)
			assert.Equal(t, tc.wantAPIs, apis)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
