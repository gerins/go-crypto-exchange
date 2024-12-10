package redis

import (
	"context"

	"github.com/gerins/log"
	"github.com/redis/go-redis/v9"

	"core-engine/config"
)

func Init(cfg config.Cache) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.Database,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("failed connecting to redis server")
	}

	return redisClient
}
