package ioc

import (
	"TestCopilot/TestEngine/config"
	redisv9 "github.com/redis/go-redis/v9"
)

func InitRedis() redisv9.Cmdable {
	redisClient := redisv9.NewClient(&redisv9.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}

//func InitRLockClient(cmd redis.Cmdable) *rlock.Client {
//	return rlock.NewClient(cmd)
//}
