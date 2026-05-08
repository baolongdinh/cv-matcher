package utils

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(opt)

	return RedisClient.Ping(context.Background()).Err()
}
