package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	"time"
)

type RankingLocalCache struct {
	// 使用泛型封装
	topN       *atomicx.Value[[]domain.Note]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		topN: atomicx.NewValue[[]domain.Note](),
		ddl:  atomicx.NewValueOf(time.Now()),
		// 用不过期，或者非常长，对齐都 redis 的过期时间
		expiration: time.Minute * 10,
	}
}

func (r RankingLocalCache) Set(ctx context.Context, notes []domain.Note) error {
	// 也可以按照 id => Note 缓存
	r.topN.Store(notes)
	ddl := time.Now().Add(r.expiration)
	r.ddl.Store(ddl)
	return nil
}

func (r RankingLocalCache) Get(ctx context.Context) ([]domain.Note, error) {
	ddl := r.ddl.Load()
	notes := r.topN.Load()
	if len(notes) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存未命中")
	}
	return notes, nil
}

func (r RankingLocalCache) Preload(ctx context.Context) {
	// 处理预加载
	panic("implement me")
}

func (r RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Note, error) {
	// 不检查过期时间，直接返回
	notes := r.topN.Load()
	return notes, nil
}

type item struct {
	notes []domain.Note
	ddl   time.Time
}
