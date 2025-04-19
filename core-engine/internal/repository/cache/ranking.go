package cache

import (
	"TestCopilot/TestEngine/internal/domain"
	"context"
	"encoding/json"
	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, notes []domain.Note) error
	Get(ctx context.Context) ([]domain.Note, error)
}

type RankingRedisCache struct {
	client redisv9.Cmdable
	key    string
}

func NewRankingRedisCache(client redisv9.Cmdable) *RankingRedisCache {
	return &RankingRedisCache{
		client: client,
		key:    "ranking",
	}
}

func (r *RankingRedisCache) Set(ctx context.Context, notes []domain.Note) error {
	for i := 0; i < len(notes); i++ {
		notes[i].Content = ""
	}
	val, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	// 过期时间要稍微长一点，最好是超过计算热榜的时间（包含重试在内的时间）
	return r.client.Set(ctx, r.key, val, time.Minute*10).Err()
}

func (r *RankingRedisCache) Get(ctx context.Context) ([]domain.Note, error) {
	data, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Note
	err = json.Unmarshal(data, &res)
	return res, err
}
