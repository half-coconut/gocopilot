package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/half-coconut/gocopilot/core-engine/internal/domain"
	redisv9 "github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Set(ctx context.Context, user domain.User) error
	Get(ctx context.Context, id int64) (domain.User, error)
	Delete(ctx context.Context, id int64) error
}

func NewUserCache(cmd redisv9.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 5, // 5min 过期时间
	}
}

type RedisUserCache struct {
	cmd        redisv9.Cmdable
	expiration time.Duration
}

func (cache *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := cache.key(user.Id)
	return cache.cmd.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	//ctx = context.WithValue(ctx, "biz", "user")
	//ctx = context.WithValue(ctx, "pattern", "user:info:%d")
	key := cache.key(id)
	val, err := cache.cmd.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err

}

func (cache *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return cache.cmd.Del(ctx, cache.key(id)).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
