package repository

import (
	"TestCopilot/TestEngine/internal/domain"
	"TestCopilot/TestEngine/internal/repository/cache"
	"context"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, notes []domain.Note) error
	GetTopN(ctx context.Context, notes []domain.Note) ([]domain.Note, error)
}
type CacheRankingRepository struct {
	// 可读性更好，对测试不友好
	redis *cache.RankingRedisCache
	local *cache.RankingLocalCache
}

func NewCacheRankingRepository(
	local *cache.RankingLocalCache,
	redis *cache.RankingRedisCache) RankingRepository {
	return &CacheRankingRepository{
		local: local,
		redis: redis}
}

func (c *CacheRankingRepository) ReplaceTopN(ctx context.Context, notes []domain.Note) error {
	// 先操作本地缓存
	_ = c.local.Set(ctx, notes)

	return c.redis.Set(ctx, notes)
}

func (c *CacheRankingRepository) GetTopN(ctx context.Context, notes []domain.Note) ([]domain.Note, error) {
	data, err := c.local.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		_ = c.local.Set(ctx, data)
	} else {
		return c.local.ForceGet(ctx)
	}
	return data, err
}
