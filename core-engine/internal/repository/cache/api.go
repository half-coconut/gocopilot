package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

type APICache interface {
	Set(ctx context.Context, api domain.API) error
	Get(ctx context.Context, id int64) (domain.API, error)
	Delete(ctx context.Context, id int64) error
}

func NewAPICache(cmd redisv9.Cmdable) APICache {
	return &RedisAPICache{
		cmd:        cmd,
		expiration: time.Minute * 5, // 5min 过期时间
	}
}

type RedisAPICache struct {
	cmd        redisv9.Cmdable
	expiration time.Duration
}

func (cache *RedisAPICache) Set(ctx context.Context, api domain.API) error {
	val, err := json.Marshal(api)
	if err != nil {
		return err
	}
	key := cache.key(api.Id)
	return cache.cmd.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisAPICache) Get(ctx context.Context, id int64) (domain.API, error) {
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.API{}, err
	}
	var api domain.API
	err = json.Unmarshal(val, &api)
	return api, err

}

func (cache *RedisAPICache) Delete(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisAPICache) key(id int64) string {
	return fmt.Sprintf("api:info:%d", id)
}
