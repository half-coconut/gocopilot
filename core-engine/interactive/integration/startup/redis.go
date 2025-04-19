package startup

import (
	"context"
	redisv9 "github.com/redis/go-redis/v9"
)

var redisClient redisv9.Cmdable

func InitRedis() redisv9.Cmdable {
	if redisClient == nil {
		redisClient = redisv9.NewClient(&redisv9.Options{
			Addr: "localhost:6379",
		})

		for err := redisClient.Ping(context.Background()).Err(); err != nil; {
			panic(err)
		}
	}
	return redisClient
}
