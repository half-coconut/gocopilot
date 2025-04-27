package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

type TaskCache interface {
	Set(ctx context.Context, task domain.Task) error
	Get(ctx context.Context, id int64) (domain.Task, error)
	Delete(ctx context.Context, id int64) error
}

func NewTaskCache(cmd redisv9.Cmdable) TaskCache {
	return &RedisTaskCache{
		cmd:        cmd,
		expiration: time.Minute * 5, // 5min 过期时间
	}
}

type RedisTaskCache struct {
	cmd        redisv9.Cmdable
	expiration time.Duration
}

func (cache *RedisTaskCache) Set(ctx context.Context, task domain.Task) error {
	val, err := json.Marshal(task)
	if err != nil {
		return err
	}
	key := cache.key(task.Id)
	return cache.cmd.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisTaskCache) Get(ctx context.Context, id int64) (domain.Task, error) {
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Task{}, err
	}
	var task domain.Task
	err = json.Unmarshal(val, &task)
	return task, err

}

func (cache *RedisTaskCache) Delete(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisTaskCache) key(id int64) string {
	return fmt.Sprintf("task:info:%d", id)
}
